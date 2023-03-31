import React, { PureComponent } from 'react'
import { Drawer, Table, Tag } from 'antd'
import { Trans } from '@lingui/macro'
import PropTypes from 'prop-types'
import dayjs from 'dayjs'
import { getTaskStatusLabel, getTaskStatusColor } from 'utils/constant'

class TaskList extends PureComponent {
  render() {
    const { onOk, dataSource, viewAlgoHandler, ...taskListProps } = this.props
    const columns = [
      {
        title: <Trans>Index</Trans>,
        dataIndex: 'index',
        key: 'index',
        width: '10%',
        fixed: 'left',
      },
      {
        title: <Trans>Running Time</Trans>,
        dataIndex: 'created_at',
        key: 'created_at',
        render: (text, _) => {
          return <span>{dayjs(text).format('YYYY-MM-DD HH:mm:ss')}</span>
        },
      },
      {
        title: <Trans>Running User</Trans>,
        dataIndex: 'username',
        key: 'username',
      },
      {
        title: <Trans>Running Status</Trans>,
        dataIndex: 'status',
        key: 'status',
        render: (text, _) => (
          <Tag color={getTaskStatusColor(text)}>{getTaskStatusLabel(text)}</Tag>
        ),
      },
      {
        title: <Trans>Operation</Trans>,
        key: 'operation',
        fixed: 'right',
        width: '20%',
        render: (text, record) => (
          <a
            onClick={() => {
              viewAlgoHandler(record.task_id)
            }}
          >
            查看运行结果
          </a>
        ),
      },
    ]

    return (
      <Drawer {...taskListProps}>
        <Table columns={columns} simple dataSource={dataSource} />
      </Drawer>
    )
  }
}

TaskList.propTypes = {
  item: PropTypes.object,
  onOk: PropTypes.func,
}

export default TaskList
