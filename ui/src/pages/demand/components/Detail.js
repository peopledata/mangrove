import React, { PureComponent } from 'react'
import { Drawer, Row, Col, Tabs, Table, Tag, message } from 'antd'
import { t, Trans } from '@lingui/macro'
import PropTypes from 'prop-types'
import styles from './Detail.less'
import dayjs from 'dayjs'
import { CopyOutlined } from '@ant-design/icons'
import {
  getDemandCategoryLabel,
  getTaskStatusLabel,
  getTaskStatusColor,
} from 'utils/constant'

class DemandDetail extends PureComponent {
  onTabChange = (key) => {
    console.log(key)
  }

  copyText = (text) => {
    navigator.clipboard.writeText(text)
    setTimeout(() => message.info(t`Copy Successfully`), 200)
  }

  viewDetailHandler = (content) => {
    const { viewDetailHandler } = this.props
    viewDetailHandler(content)
  }

  renderDemandInfo = () => {
    const { demand } = this.props
    const { demandDetail } = demand
    if (!demandDetail) {
      return <></>
    }
    return (
      <div className={styles.detail}>
        <p className={styles.headerLabel}>{t`Basic Info`}</p>
        <div className={styles.topContent}>
          <Row gutter={24}>
            <Col span={12}>
              <div className={styles.item}>
                <label>{t`Demand Name`}：</label>
                <span>{demandDetail.name}</span>
              </div>
            </Col>
            <Col span={12}>
              <div className={styles.item}>
                <label>{t`Demand Valid Time`}：</label>
                <span>
                  {dayjs(demandDetail.valid_at).format('YYYY-MM-DD HH:mm:ss')}
                </span>
              </div>
            </Col>
          </Row>

          <Row gutter={24}>
            <Col span={12}>
              <div className={styles.item}>
                <label>{t`Data Category`}：</label>
                <span>{getDemandCategoryLabel(demandDetail.category)}</span>
              </div>
            </Col>
            <Col span={12}>
              <div className={styles.item}>
                <label>{t`Data Content`}：</label>
                <span>{demandDetail.content}</span>
              </div>
            </Col>
          </Row>

          <Row gutter={24}>
            <Col span={12}>
              <div className={styles.item}>
                <label>{t`Needs Users`}：</label>
                <span>{demandDetail.need_users}</span>
              </div>
            </Col>
            <Col span={12}>
              <div className={styles.item}>
                <label>{t`Existing Users`}：</label>
                <span>{demandDetail.existing_users}</span>
              </div>
            </Col>
          </Row>

          <Row gutter={24}>
            <Col span={12}>
              <div className={styles.item}>
                <label>{t`Data Available Times`}：</label>
                <span>{demandDetail.available_times}</span>
              </div>
            </Col>
            <Col span={12}>
              <div className={styles.item}>
                <label>{t`Used Times`}：</label>
                <span>{demandDetail.use_times}</span>
              </div>
            </Col>
          </Row>

          <Row gutter={24}>
            <Col span={12}>
              <div className={styles.item}>
                <label>{t`Create User`}：</label>
                <span>xxx</span>
              </div>
            </Col>
            <Col span={12}>
              <div className={styles.item}>
                <label>{t`CreateTime`}：</label>
                <span>
                  {dayjs(demandDetail.created_at).format('YYYY-MM-DD HH:mm:ss')}
                </span>
              </div>
            </Col>
          </Row>

          <Row gutter={24}>
            <Col span={12}>
              <div className={styles.item}>
                <label>{t`Blockchain ID`}：</label>
                <span>{demandDetail.contract_symbol}</span>
              </div>
            </Col>
            <Col span={12}>
              <div className={styles.item}>
                <label>{t`Contract Address`}：</label>
                {demandDetail.contract_address !== '' ? (
                  <>
                    <span>
                      {demandDetail.contract_address.slice(0, 15)}......
                    </span>
                    <CopyOutlined
                      onClick={() => {
                        this.copyText(demandDetail.contract_address)
                      }}
                    />
                  </>
                ) : (
                  <span>{t`No Data`}</span>
                )}
              </div>
            </Col>
          </Row>

          <Row gutter={24}>
            <Col span={24}>
              <div className={styles.item}>
                <label>{t`Data Usage`}：</label>
                <a
                  onClick={() => {
                    this.viewDetailHandler(demandDetail.purpose)
                  }}
                >{t`View Detail`}</a>
              </div>
            </Col>
          </Row>
        </div>

        <p
          style={{ marginTop: '100px' }}
          className={styles.headerLabel}
        >{t`Algorithms and protocols`}</p>
        <div>
          <div className={styles.item}>
            <label>{t`Algorithm file image address`}：</label>
            <span>{demandDetail.Algorithm}</span>
          </div>
          <div className={styles.item}>
            <label>{t`Agreement Content`}：</label>
            <a
              onClick={() => {
                this.viewDetailHandler(demandDetail.agreement)
              }}
            >{t`View Detail`}</a>
          </div>
        </div>
      </div>
    )
  }

  tableChange = (page) => {
    const { demand, recordsTableChange } = this.props
    const { demandDetail } = demand
    recordsTableChange(page, demandDetail.demand_id)
  }

  renderContactRecords = () => {
    const { demand } = this.props
    const { contractRecords, contractRecordsPager } = demand
    const columns = [
      {
        title: <Trans>User DID</Trans>,
        dataIndex: 'did',
        key: 'did',
        fixed: 'left',
        render: (text, _) => <span>{text}</span>,
      },
      {
        title: <Trans>Sign Time</Trans>,
        dataIndex: 'sign_time',
        key: 'sign_time',
        render: (text, _) => {
          return <span>{dayjs(text).format('YYYY-MM-DD HH:mm:ss')}</span>
        },
      },
    ]
    return (
      <Table
        contractRecordsPager
        pagination={{
          ...contractRecordsPager,
          showSizeChanger: true,
          showTotal: (total) => t`Total ${total} Items`,
        }}
        dataSource={contractRecords}
        columns={columns}
        simple
        rowKey={(record) => record.id}
        onChange={this.tableChange}
      />
    )
  }

  renderTasks = () => {
    const { demand } = this.props
    const { taskList } = demand
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
        width: '15%',
      },
    ]
    return <Table columns={columns} simple dataSource={taskList} />
  }

  render() {
    const { onOk, dataSource, ...drawerProps } = this.props
    const items = [
      {
        key: '1',
        label: <Trans>Demand Info</Trans>,
        children: this.renderDemandInfo(),
      },
      {
        key: '2',
        label: <Trans>Contracted User</Trans>,
        children: this.renderContactRecords(),
      },
      {
        key: '3',
        label: <Trans>Data Operation Record</Trans>,
        children: this.renderTasks(),
      },
    ]
    return (
      <Drawer {...drawerProps}>
        <Tabs defaultActiveKey="1" items={items} onChange={this.onTabChange} />
      </Drawer>
    )
  }
}

DemandDetail.propTypes = {
  demand: PropTypes.object,
}

export default DemandDetail
