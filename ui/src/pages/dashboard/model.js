import { parse } from 'qs'
import modelExtend from 'dva-model-extend'
import api from 'api'
import { pathToRegexp } from 'path-to-regexp'
import { model } from 'utils/model'

const { queryDashboard } = api

export default modelExtend(model, {
  namespace: 'dashboard',
  state: {
    sales: [],
  },
  subscriptions: {
    setup({ dispatch, history }) {
      history.listen(({ pathname }) => {
        if (
          pathToRegexp('/dashboard').exec(pathname) ||
          pathToRegexp('/').exec(pathname)
        ) {
          dispatch({ type: 'query' })
        }
      })
    },
  },
  effects: {
    *query({ payload }, { call, put }) {
      const response = yield call(queryDashboard, parse(payload))
      yield put({
        type: 'updateState',
        payload: response.data,
      })
    },
  },
})
