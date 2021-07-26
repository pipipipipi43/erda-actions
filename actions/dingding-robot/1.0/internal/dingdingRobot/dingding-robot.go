package dingdingRobot

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/blinkbean/dingtalk"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/dingding-robot/1.0/internal/conf"
	"github.com/erda-project/erda/pkg/strutil"
)

// phRe 占位符正则表达式:
//   ${{ configs.key }}
//   ${{ dirs.preTaskName.fileName }}
//   ${{ outputs.preTaskName.key }}
//   ${{ params.key }}
//   ${{ (echo hello world) }}
var PhRe = regexp.MustCompile(`\${{[ ]{1}([^{}\s]+)[ ]{1}}}`) // [ ]{1} 强调前后均有且仅有一个空格

// loosePhRe 宽松的正则表达式:
var LoosePhRe = regexp.MustCompile(`\${{[^{}]+}}`)

func handleAPIs() error {
	params := map[string]string{
		"env.DiceOpenapiAddr":      conf.DiceOpenapiAddr(),
		"env.Pipeline_Task_Id":     strconv.FormatUint(conf.TaskId(), 10),
		"env.MetaFile":             conf.MetaFile(),
		"env.WorkDir":              conf.WorkDir(),
		"env.DiceVersion":          conf.DiceVersion(),
		"env.PipelineId":           strconv.FormatUint(conf.PipelineId(), 10),
		"env.OrgId":                strconv.FormatUint(conf.OrgId(), 10),
		"env.TaskId":               strconv.FormatUint(conf.TaskId(), 10),
		"env.ProjectId":            strconv.FormatUint(conf.ProjectId(), 10),
		"env.SponsorId":            conf.SponsorId(),
		"env.CommitId":             conf.CommitId(),
		"env.GittarUsername":       conf.GittarUsername(),
		"env.GittarPassword":       conf.GittarPassword(),
		"env.ProjectName":          conf.ProjectName(),
		"env.BranchName":           conf.BranchName(),
		"env.DiceClusterName":      conf.DiceClusterName(),
		"env.DiceOpenapiToken":     conf.DiceOpenapiToken(),
		"env.DiceOpenapiPublicUrl": conf.DiceOpenapiPublicUrl(),
	}
	msg, err := Eval(conf.Msg(), params)
	if err != nil {
		return err
	}

	sendMsg(msg)
	logrus.Info("发送消息成功")
	return nil
}

func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, s string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(s), -1) {
		var groups []string
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, s[v[i]:v[i+1]])
		}

		result += s[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + s[lastIndex:]
}

// FindInvalidPlaceholders 找到表达式中不合规范的占位符
func FindInvalidPlaceholders(exprStr string) []string {
	// 合法的占位符列表
	validPhs := PhRe.FindAllString(exprStr, -1)
	// 宽松的占位符列表
	loosePhs := LoosePhRe.FindAllString(exprStr, -1)

	// 不在合法占位符列表中的即为非法占位符
	var invalidPhs []string
	for _, loose := range loosePhs {
		if !strutil.Exist(validPhs, loose) {
			invalidPhs = append(invalidPhs, loose)
		}
	}
	return invalidPhs
}

func Eval(exprStr string, placeholderParams map[string]string) (string, error) {
	// 校验表达式
	invalidPhs := FindInvalidPlaceholders(exprStr)
	if len(invalidPhs) > 0 {
		return "", fmt.Errorf("invalid expression, found invalid placeholders: %s (must match: %s)", strings.Join(invalidPhs, ", "), PhRe.String())
	}

	// 递归渲染表达式中的占位符
	// 直到渲染完毕或找到不存在的占位符
	for {
		var notFoundPlaceholders []string
		exprStr = strutil.ReplaceAllStringSubmatchFunc(PhRe, exprStr, func(subs []string) string {
			ph := subs[0]    // ${{ configs.key }}
			inner := subs[1] // configs.key

			v, ok := placeholderParams[inner]
			if ok {
				return v
			}
			notFoundPlaceholders = append(notFoundPlaceholders, ph)
			return exprStr
		})
		if len(notFoundPlaceholders) > 0 {
			return "", fmt.Errorf("invalid expression, not found placeholders: %s", strings.Join(notFoundPlaceholders, ", "))
		}
		// 没有需要替换的占位符，则退出渲染
		if !PhRe.MatchString(exprStr) {
			break
		}
	}

	return exprStr, nil
}

func sendMsg(msg string) error {
	var dingToken = conf.RobotToken()
	cli := dingtalk.InitDingTalk(dingToken, conf.KeyWord())
	err := cli.SendTextMessage(msg)
	if err != nil {
		logrus.Info("发送消息失败")
		return err
	}
	return nil
}
