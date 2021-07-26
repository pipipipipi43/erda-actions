通过Action发送消息
在Msg中填写需要发送的消息,可以使用占位符

#支持以下表达式

${{ configs.key }}
${{ dirs.preTaskName.fileName }}
${{ outputs.preTaskName.key }}
${{ params.key }}
${{ (echo hello world) }}
${{ env.key }}

#例子
      - dingding-robot:
          alias: dingding-robot
          description: 给指定的钉钉群发送消息
          version: "1.0"
          params:
            Keyword: .
            Msg: ${{ env.ProjectId }}/${{ env.PipelineId }}
            Token:
              - 2d6824abbc5c3c495b84c9d87cec649f7f559281a86c48fa4d2836c5d449fc5f
    