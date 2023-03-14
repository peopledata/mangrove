import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { Table, Modal, Tag } from 'antd'
import { DropOption } from 'components'
import dayjs from 'dayjs'
import { t } from '@lingui/macro'
import { Trans } from '@lingui/macro'
import { Link } from 'umi'
import {
  getDemandStatusLabel,
  getDemandStatusColor,
  getDemandCategoryLabel,
} from 'utils/constant'
import styles from './List.less'
import { DEMAND_STATUS } from 'utils/constant'

const { confirm } = Modal

class List extends PureComponent {
  handleMenuClick = (record, e) => {
    const {
      onDeleteItem,
      onEditItem,
      onExecuteItem,
      onPublishItem,
      onCloseItem,
      onRecordItem,
    } = this.props
    if (e.key === '1') {
      // 发布
      confirm({
        title: t`Are you sure publish this demand?`,
        onOk() {
          onPublishItem(record.demand_id)
        },
      })
    } else if (e.key === '2') {
      // 更新
      onEditItem(record.demand_id)
    } else if (e.key === '3') {
      // 删除
      confirm({
        title: t`Are you sure delete these demands?`,
        onOk() {
          onDeleteItem(record.demand_id)
        },
      })
    } else if (e.key === '4') {
      //   执行计算
      confirm({
        title: t`Are you sure execute data?`,
        onOk() {
          onExecuteItem(record.demand_id)
        },
      })
    } else if (e.key === '5') {
      //   下架
      confirm({
        title: t`Are you sure close this demand?`,
        onOk() {
          onCloseItem(record.demand_id)
        },
      })
    } else if (e.key === '6') {
      //   查看运行记录
      onRecordItem(record.demand_id)
    }
  }

  handlerItemClick = (record, e) => {
    const { onDetailItem } = this.props
    onDetailItem(record.demand_id)
  }

  operationOptions = (record) => {
    if (record.status === DEMAND_STATUS.PUBLISHED) {
      return [
        { key: '5', name: t`Close` },
        { key: '4', name: t`Execute` },
        { key: '6', name: t`Record` },
      ]
    } else if (record.status === DEMAND_STATUS.COMPLETED) {
      return [
        { key: '4', name: t`Execute` },
        { key: '6', name: t`Record` },
      ]
    } else if (record.status === DEMAND_STATUS.CLOSED) {
      return [{ key: '3', name: t`Delete` }]
    }
    // INIT状态
    return [
      { key: '1', name: t`Publish` },
      { key: '2', name: t`Update` },
      { key: '3', name: t`Delete` },
    ]
  }

  render() {
    const { ...tableProps } = this.props

    const columns = [
      {
        title: <Trans>Demand Name</Trans>,
        dataIndex: 'name',
        key: 'name',
        width: '10%',
        fixed: 'left',
        render: (text, record) => (
          <div
            onClick={(e) => this.handlerItemClick(record, e)}
            className={styles.name}
          >
            {text}
          </div>
        ),
      },
      {
        title: <Trans>Status</Trans>,
        dataIndex: 'status',
        key: 'status',
        render: (text, _) => (
          <Tag color={getDemandStatusColor(text)}>
            {getDemandStatusLabel(text)}
          </Tag>
        ),
      },
      {
        title: <Trans>Category</Trans>,
        dataIndex: 'category',
        key: 'category',
        render: (text, _) => <span>{getDemandCategoryLabel(text)}</span>,
      },
      {
        title: <Trans>Needs Users</Trans>,
        dataIndex: 'need_users',
        key: 'need_users',
      },
      {
        title: <Trans>Existing Users</Trans>,
        dataIndex: 'existing_users',
        key: 'existing_users',
      },
      {
        title: <Trans>Data Available Times</Trans>,
        dataIndex: 'available_times',
        key: 'available_times',
      },
      {
        title: <Trans>Used Times</Trans>,
        dataIndex: 'use_times',
        key: 'use_times',
      },
      {
        title: <Trans>ValidTime</Trans>,
        dataIndex: 'valid_at',
        key: 'valid_at',
        width: '14%',
        render: (text, _) => {
          return <span>{dayjs(text).format('YYYY-MM-DD HH:mm:ss')}</span>
        },
      },
      {
        title: <Trans>CreateTime</Trans>,
        dataIndex: 'created_at',
        key: 'created_at',
        width: '14%',
        render: (text, _) => {
          return <span>{dayjs(text).format('YYYY-MM-DD HH:mm:ss')}</span>
        },
      },
      {
        title: <Trans>Operation</Trans>,
        key: 'operation',
        fixed: 'right',
        width: '8%',
        render: (text, record) => {
          return record.status === DEMAND_STATUS.PUBLISHING ? (
            <></>
          ) : (
            <DropOption
              onMenuClick={(e) => this.handleMenuClick(record, e)}
              menuOptions={this.operationOptions(record)}
            />
          )
        },
      },
    ]

    return (
      <Table
        {...tableProps}
        pagination={{
          ...tableProps.pagination,
          showTotal: (total) => t`Total ${total} Items`,
        }}
        className={styles.table}
        // bordered
        columns={columns}
        simple
        rowKey={(record) => record.id}
      />
    )
  }
}

List.propTypes = {
  onDeleteItem: PropTypes.func,
  onEditItem: PropTypes.func,
  onDetailItem: PropTypes.func,
  onPublishItem: PropTypes.func,
  onExecuteItem: PropTypes.func,
  onCloseItem: PropTypes.func,
  location: PropTypes.object,
}

export default List
