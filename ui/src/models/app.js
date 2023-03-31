/* global window */
import { history } from 'umi'
import { stringify } from 'qs'
import store from 'store'
import { ROLE_TYPE } from 'utils/constant'
import { queryLayout } from 'utils'
import { CANCEL_REQUEST_MESSAGE } from 'utils/constant'
import api from 'api'
import config from 'config'
import { pathToRegexp } from 'path-to-regexp'
import { CODE_SUCCESS } from '../utils/constant'

const { queryRouteList, logoutUser, queryUserInfo } = api

const goDemand = () => {
  if (pathToRegexp(['/', '/login']).exec(window.location.pathname)) {
    history.push({
      pathname: '/demand',
    })
  }
}

export default {
  namespace: 'app',
  state: {
    routeList: [],
    locationPathname: '',
    locationQuery: {},
    theme: store.get('theme') || 'light',
    collapsed: store.get('collapsed') || false,
  },
  subscriptions: {
    setup({ dispatch }) {
      dispatch({ type: 'query' })
    },
    setupHistory({ dispatch, history }) {
      history.listen((location) => {
        dispatch({
          type: 'updateState',
          payload: {
            locationPathname: location.pathname,
            locationQuery: location.query,
          },
        })
      })
    },

    setupRequestCancel({ history }) {
      history.listen(() => {
        const { cancelRequest = new Map() } = window
        cancelRequest.forEach((value, key) => {
          if (value.pathname !== window.location.pathname) {
            value.cancel(CANCEL_REQUEST_MESSAGE)
            cancelRequest.delete(key)
          }
        })
      })
    },
  },

  effects: {
    *query({ payload }, { call, put, select }) {
      // store isInit to prevent query trigger by refresh
      const isInit = store.get('isInit')
      console.log('query.isInit=', isInit)
      if (isInit) {
        goDemand()
        return
      }
      const { locationPathname } = yield select((_) => _.app)
      const { success, code, data } = yield call(queryUserInfo, payload)
      const user = data
      // true 1008 undefined
      if (success && code === CODE_SUCCESS && user) {
        const { data } = yield call(queryRouteList)
        const { permissions } = user
        let routeList = data
        if (
          permissions.role === ROLE_TYPE.ADMIN ||
          permissions.role === ROLE_TYPE.DEVELOPER
        ) {
          permissions.visit = data.map((item) => item.id)
        } else {
          routeList = data.filter((item) => {
            const cases = [
              permissions.visit.includes(item.id),
              item.mpid
                ? permissions.visit.includes(item.mpid) || item.mpid === '-1'
                : true,
              item.bpid ? permissions.visit.includes(item.bpid) : true,
            ]
            return cases.every((_) => _)
          })
        }
        store.set('routeList', routeList)
        store.set('permissions', permissions)
        store.set('user', user)
        store.set('isInit', true)
        goDemand()
      } else if (queryLayout(config.layouts, locationPathname) !== 'public') {
        console.log('====push2login====')
        history.push({
          pathname: '/login',
          search: stringify({
            from: locationPathname,
          }),
        })
      }
    },

    *signOut({ payload }, { call, put, select }) {
      store.remove('routeList')
      store.remove('permissions')
      store.remove('user')
      store.remove('isInit')
      store.remove('access_token')
      store.remove('refresh_token')
      const { locationPathname } = yield select((_) => _.app)
      history.push({
        pathname: '/login',
        search: stringify({
          from: locationPathname,
        }),
      })
    },
  },

  reducers: {
    updateState(state, { payload }) {
      return {
        ...state,
        ...payload,
      }
    },

    handleThemeChange(state, { payload }) {
      store.set('theme', payload)
      state.theme = payload
    },

    handleCollapseChange(state, { payload }) {
      store.set('collapsed', payload)
      state.collapsed = payload
    },
  },
}
