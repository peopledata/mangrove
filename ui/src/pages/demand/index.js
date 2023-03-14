import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { Page } from 'components'
import { history } from 'umi'
import { connect } from 'umi'
import { stringify } from 'qs'
import List from './components/List'
import Filter from './components/Filter'
import Modal from './components/Modal'
import Drawer from './components/Drawer'
import Detail from './components/Detail'
import { t } from '@lingui/macro'
import { Button, Col, Popconfirm, Row } from 'antd'

@connect(({ demand, loading }) => ({ demand, loading }))
class Demand extends PureComponent {
  handleRefresh = (newQuery) => {
    const { location } = this.props
    const { query, pathname } = location

    history.push({
      pathname,
      search: stringify(
        {
          ...query,
          ...newQuery,
        },
        { arrayFormat: 'repeat' }
      ),
    })
  }

  get filterProps() {
    const { location, dispatch } = this.props
    const { query } = location

    return {
      filter: {
        ...query,
      },
      onFilterChange: (value) => {
        this.handleRefresh({
          ...value,
        })
      },
      onAdd() {
        dispatch({
          type: 'demand/showModal',
          payload: {
            modalType: 'create',
          },
        })
      },
      onRefresh() {
        dispatch({
          type: 'demand/query',
          payload: {},
        })
      },
    }
  }

  get listProps() {
    const { dispatch, demand, loading } = this.props
    const { list, pagination, contractRecordsPager } = demand

    return {
      dataSource: list,
      loading: loading.effects['demand/query'],
      pagination,
      onChange: (page) => {
        dispatch({
          type: 'demand/query',
          payload: {
            page: page.current,
            pageSize: page.pageSize,
          },
        })
      },
      onDeleteItem: (demandId) => {
        dispatch({
          type: 'demand/delete',
          payload: demandId,
        }).then(() => {
          dispatch({
            type: 'demand/query',
            payload: {
              page:
                list.length === 1 && pagination.current > 1
                  ? pagination.current - 1
                  : pagination.current,
              pageSize: pagination.pageSize,
            },
          })
        })
      },
      onEditItem(demandId) {
        dispatch({
          type: 'demand/detail',
          payload: {
            modalType: 'update',
            id: demandId,
          },
        })
      },
      onPublishItem(demandId) {
        dispatch({
          type: 'demand/publish',
          payload: {
            id: demandId,
          },
        })
      },
      onExecuteItem(demandId) {
        console.log('onExecuteItem.demandId=', demandId)
      },
      onCloseItem(demandId) {
        console.log('onCloseItem.demandId=', demandId)
      },
      onRecordItem(demandId) {
        dispatch({
          type: 'demand/queryTask',
          payload: {
            id: demandId,
          },
        })
      },
      onDetailItem(demandId) {
        dispatch({
          type: 'demand/queryDetail',
          payload: {
            id: demandId,
          },
        })
        dispatch({
          type: 'demand/queryContractRecords',
          payload: {
            id: demandId,
            page: contractRecordsPager.current,
            pageSize: contractRecordsPager.pageSize,
          },
        })
      },
    }
  }

  get modalProps() {
    const { dispatch, demand, loading } = this.props
    const { currentItem, modalOpen, modalType } = demand

    return {
      item: modalType === 'create' ? {} : currentItem,
      open: modalOpen,
      destroyOnClose: true,
      maskClosable: false,
      confirmLoading: loading.effects[`demand/${modalType}`],
      title: `${modalType === 'create' ? t`Create Demand` : t`Update Demand`}`,
      centered: true,
      width: 680,
      bodyStyle: { overflowX: 'hidden', overflowY: 'scroll', height: '77vh' },
      onOk: (data) => {
        dispatch({
          type: `demand/${modalType}`,
          payload: data,
        }).then(() => {
          this.handleRefresh()
        })
      },
      onCancel() {
        dispatch({
          type: 'demand/hideModal',
        })
      },
    }
  }

  get drawerProps() {
    const { dispatch, demand } = this.props
    const { drawerOpen } = demand
    return {
      title: t`Operation Record`,
      size: 'large',
      placement: 'right',
      open: drawerOpen,
      onClose: () => {
        dispatch({
          type: 'demand/hideDrawer',
        })
      },
    }
  }

  get detailProps() {
    const { dispatch, demand } = this.props
    const { detailOpen } = demand
    return {
      title: t`Demand Detail`,
      size: 'large',
      placement: 'right',
      open: detailOpen,
      onClose: () => {
        dispatch({
          type: 'demand/hideDetail',
        })
      },
    }
  }

  recordsTableChange = (page, demandId) => {
    const { dispatch } = this.props
    dispatch({
      type: 'demand/queryContractRecords',
      payload: {
        id: demandId,
        page: page.current,
        pageSize: page.pageSize,
      },
    })
  }

  render() {
    const { demand } = this.props
    const { taskList } = demand

    return (
      <Page inner>
        <Filter {...this.filterProps} />
        <List {...this.listProps} />
        <Modal {...this.modalProps} />
        <Drawer {...this.drawerProps} dataSource={taskList} />
        <Detail
          {...this.detailProps}
          demand={demand}
          recordsTableChange={this.recordsTableChange}
        />
      </Page>
    )
  }
}

Demand.propTypes = {
  demand: PropTypes.object,
  location: PropTypes.object,
  dispatch: PropTypes.func,
  loading: PropTypes.object,
}

export default Demand
