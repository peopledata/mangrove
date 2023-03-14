import modelExtend from 'dva-model-extend'
import { pathToRegexp } from 'path-to-regexp'
import dayjs from 'dayjs'
import api from 'api'
import { pageModel } from 'utils/model'
import { CODE_SUCCESS } from '../../utils/constant'

const {
  queryDemandList,
  queryTaskList,
  queryContractRecords,
  queryDemandInfo,
  queryDemandDetail,
  createDemand,
  updateDemand,
  publishDemand,
  removeDemand,
  removeDemandList,
} = api

export default modelExtend(pageModel, {
  namespace: 'demand',

  state: {
    taskList: [],
    contractRecords: [],
    contractRecordsPager: {
      showSizeChanger: true,
      showQuickJumper: true,
      current: 1,
      total: 0,
      pageSize: 10,
    },
    currentItem: {},
    modalOpen: false,
    modalType: 'create',
    drawerOpen: false,
    detailOpen: false,
    drawerShowDetail: false,
    showDetailContent: '',
    demandDetail: null,
    selectedRowKeys: [],
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen((location) => {
        if (pathToRegexp('/demand').exec(location.pathname)) {
          // const payload = location.query || { page: 1, pageSize: 10 }
          const payload = { page: 1, pageSize: 10 }
          dispatch({
            type: 'query',
            payload,
          })
        }
      })
    },
  },

  effects: {
    *query({ payload = {} }, { call, put }) {
      const page = Number(payload.page) || 10
      const pageSize = Number(payload.pageSize) || 10
      const q = payload.q || ''
      const response = yield call(queryDemandList, {
        query: `page=${page}&pageSize=${pageSize}&q=${q}`,
      })
      if (response.data) {
        const { demands, total } = response.data
        console.log('demands=', demands)
        yield put({
          type: 'querySuccess',
          payload: {
            list: demands,
            pagination: {
              current: Number(payload.page) || 1,
              pageSize: Number(payload.pageSize) || 10,
              total: total,
            },
          },
        })
      }
    },

    // 用于编辑更新的回填
    *detail({ payload }, { call, put }) {
      const response = yield call(queryDemandInfo, payload)
      if (response.data) {
        const demandDetail = response.data
        // datepicker 需要使用 dayjs 进行格式化
        demandDetail.valid_at = dayjs(
          demandDetail.valid_at,
          'YYYY-MM-DD HH:mm:ss'
        )
        yield put({
          type: 'showModal',
          payload: {
            modalType: 'update',
            currentItem: demandDetail,
          },
        })
      }
    },

    *delete({ payload }, { call, put, select }) {
      const data = yield call(removeDemand, { id: payload })
      const { selectedRowKeys } = yield select((_) => _.user)
      if (data.success) {
        yield put({
          type: 'updateState',
          payload: {
            selectedRowKeys: selectedRowKeys.filter((_) => _ !== payload),
          },
        })
      } else {
        throw data
      }
    },

    *multiDelete({ payload }, { call, put }) {
      const data = yield call(removeDemandList, payload)
      if (data.success) {
        yield put({ type: 'updateState', payload: { selectedRowKeys: [] } })
      } else {
        throw data
      }
    },

    *create({ payload }, { call, put }) {
      const data = yield call(createDemand, payload)
      if (data.success) {
        yield put({ type: 'hideModal' })
      } else {
        throw data
      }
    },

    *update({ payload }, { select, call, put }) {
      const id = payload.demand_id
      const newDemand = { ...payload, id }
      const data = yield call(updateDemand, newDemand)
      if (data.success && data.code === CODE_SUCCESS) {
        yield put({ type: 'hideModal' })
      } else {
        throw data
      }
    },

    *publish({ payload }, { call, put }) {
      const response = yield call(publishDemand, payload)
      if (response.success && response.code === CODE_SUCCESS) {
        // todo：获取payload
        const payload = {}
        yield put({ type: 'query', payload: payload })
      } else {
        throw response
      }
    },

    *queryTask({ payload = {} }, { call, put }) {
      const response = yield call(queryTaskList, payload)
      if (response.success && response.code === CODE_SUCCESS && response.data) {
        const { tasks, total } = response.data
        yield put({ type: 'updateState', payload: { taskList: tasks } })
      } else {
        throw response
      }
    },

    // 用于显示详情页数据
    *queryDetail({ payload = {} }, { call, put }) {
      // 用于编辑更新的回填
      const response = yield call(queryDemandDetail, payload)
      if (response.data) {
        const demandDetail = response.data
        yield put({
          type: 'showDetail',
          payload: {
            demandDetail: demandDetail,
          },
        })
      }
    },

    // 查询签约的协议用户
    *queryContractRecords({ payload = {} }, { call, put }) {
      const page = Number(payload.page) || 1
      const pageSize = Number(payload.pageSize) || 10
      payload['query'] = `page=${page}&pageSize=${pageSize}`
      const response = yield call(queryContractRecords, payload)
      if (response.success && response.code === CODE_SUCCESS && response.data) {
        const { records, total } = response.data
        yield put({
          type: 'updateState',
          payload: {
            contractRecords: records,
            contractRecordsPager: {
              current: Number(payload.page) || 1,
              pageSize: Number(payload.pageSize) || 10,
              total: total,
            },
          },
        })
      } else {
        throw response
      }
    },
  },

  reducers: {
    showModal(state, { payload }) {
      return { ...state, ...payload, modalOpen: true }
    },

    hideModal(state) {
      return { ...state, modalOpen: false }
    },

    showDrawer(state, { payload }) {
      return { ...state, ...payload, drawerOpen: true }
    },

    hideDrawer(state) {
      return { ...state, drawerOpen: false }
    },

    showDetail(state, { payload }) {
      return { ...state, ...payload, detailOpen: true }
    },

    hideDetail(state) {
      return { ...state, detailOpen: false }
    },

    drawerShowDetailOpen(state, { payload }) {
      return { ...state, ...payload, drawerShowDetail: true }
    },

    drawerShowDetailClose(state) {
      return { ...state, drawerShowDetail: false }
    },
  },
})
