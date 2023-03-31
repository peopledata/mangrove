import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { Page } from 'components'
import { history } from 'umi'
import { connect } from 'umi'
import { stringify } from 'qs'
import List from './components/List'
import Filter from './components/Filter'
import Modal from './components/Modal'
import TaskList from './components/TaskList'
import AlgoList from './components/AlgoList'
import Detail from './components/Detail'
import ShowDetail from './components/ShowDetail'
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
    const { location, dispatch, demand } = this.props
    const { list, pagination } = demand
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
          payload: {
            page:
              list && list.length === 1 && pagination.current > 1
                ? pagination.current - 1
                : pagination.current,
            pageSize: pagination.pageSize,
          },
        })
      },
      onSearch(q) {
        console.log('q=', q)
        dispatch({
          type: 'demand/query',
          payload: {
            q: q,
            page:
              list && list.length === 1 && pagination.current > 1
                ? pagination.current - 1
                : pagination.current,
            pageSize: pagination.pageSize,
          },
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
                list && list.length === 1 && pagination.current > 1
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
        dispatch({
          type: 'demand/createTask',
          payload: {
            id: demandId,
          },
        })
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
        dispatch({
          type: 'demand/showDrawer',
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
        dispatch({
          type: 'demand/queryTask',
          payload: {
            id: demandId,
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
        })
      },
      onCancel() {
        dispatch({
          type: 'demand/hideModal',
        })
      },
    }
  }

  get taskListProps() {
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

  get algoListProps() {
    const { dispatch, demand } = this.props
    const { drawerAlgoListOpen } = demand
    return {
      title: t`Operation Record`,
      size: 'large',
      placement: 'right',
      open: drawerAlgoListOpen,
      onClose: () => {
        dispatch({
          type: 'demand/drawerAlgoListClose',
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

  get showDetailProps() {
    const { dispatch, demand } = this.props
    const { drawerShowDetail } = demand
    return {
      title: t`Demand Detail`,
      size: 'large',
      placement: 'right',
      open: drawerShowDetail,
      onClose: () => {
        dispatch({
          type: 'demand/drawerShowDetailClose',
        })
        dispatch({
          type: 'demand/updateState',
          payload: {
            showDetailContent: '',
          },
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

  viewDetailHandler = (content) => {
    const { dispatch } = this.props
    dispatch({
      type: 'demand/drawerShowDetailOpen',
    })
    dispatch({
      type: 'demand/updateState',
      payload: {
        showDetailContent: content,
      },
    })
  }

  viewAlgoHandler = (taskId) => {
    const { dispatch } = this.props
    dispatch({
      type: 'demand/drawerAlgoListOpen',
    })
    dispatch({
      type: 'demand/queryAlgoRecordList',
      payload: {
        id: taskId,
      },
    })
  }

  render() {
    const { demand } = this.props
    const { taskList, showDetailContent, algoRecordList } = demand

    return (
      <Page inner>
        <Filter {...this.filterProps} />
        <List {...this.listProps} />
        <Modal {...this.modalProps} />
        <TaskList
          {...this.taskListProps}
          viewAlgoHandler={this.viewAlgoHandler}
          dataSource={taskList}
        />
        <AlgoList {...this.algoListProps} dataSource={algoRecordList} />
        <Detail
          {...this.detailProps}
          demand={demand}
          viewDetailHandler={this.viewDetailHandler}
          recordsTableChange={this.recordsTableChange}
          viewAlgoHandler={this.viewAlgoHandler}
        />
        <ShowDetail {...this.showDetailProps} content={showDetailContent} />
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
