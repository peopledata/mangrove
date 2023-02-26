import { history } from 'umi'
import api from 'api'
import store from 'store'
import { pathToRegexp } from 'path-to-regexp'
import { CODE_SUCCESS } from '../../utils/constant'

const { loginUser } = api

export default {
  namespace: 'login',

  state: {},
  // subscriptions: {
  //   setup({ dispatch, history }) {
  //     history.listen(location => {
  //       if (pathToRegexp('/login').exec(location.pathname)) {
  //       }
  //     })
  //   },
  // },
  effects: {
    *login({ payload }, { put, call, select }) {
      const data = yield call(loginUser, payload)
      const { locationQuery } = yield select((_) => _.app)
      if (data.success && data.code === CODE_SUCCESS) {
        // save jwt token to store
        console.log(data, '=======')
        const { access_token, refresh_token } = data.data
        store.set('access_token', access_token)
        store.set('refresh_token', refresh_token)
        const { from } = locationQuery
        yield put({ type: 'app/query' })
        if (!pathToRegexp('/login').exec(from)) {
          if (['', '/'].includes(from)) history.push('/dashboard')
          else history.push(from)
        } else {
          history.push('/dashboard')
        }
      } else {
        throw data
      }
    },
  },
}
