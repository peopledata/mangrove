import modelExtend from 'dva-model-extend'
import { pathToRegexp } from 'path-to-regexp'
import dayjs from 'dayjs'
import api from 'api'
import { pageModel } from 'utils/model'
import { CODE_SUCCESS } from '../../utils/constant'

const {
  queryDemandList,
  queryTaskList,
  queryDemand,
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
    currentItem: {},
    modalOpen: false,
    modalType: 'create',
    drawerOpen: false,
    selectedRowKeys: [],
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen((location) => {
        if (pathToRegexp('/demand').exec(location.pathname)) {
          const payload = location.query || { page: 1, pageSize: 10 }
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
      const response = yield call(queryDemandList, payload)
      if (response.data) {
        const { demands, total } = response.data
        yield put({
          type: 'querySuccess',
          payload: {
            list: demands,
            current: Number(payload.page) || 1,
            pageSize: Number(payload.pageSize) || 10,
            total: total,
          },
        })
      }
    },

    *detail({ payload }, { call, put }) {
      const response = yield call(queryDemand, payload)
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
        yield put({ type: 'showDrawer' })
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
  },
})
