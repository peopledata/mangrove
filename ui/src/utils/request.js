import axios from 'axios'
import { cloneDeep } from 'lodash'
import { message } from 'antd'
import { history } from 'umi'
import store from 'store'
import { CANCEL_REQUEST_MESSAGE } from 'utils/constant'
import { CODE_INVALID_TOKEN, CODE_SUCCESS } from './constant'

const { parse, compile } = require('path-to-regexp')

const { CancelToken } = axios
window.cancelRequest = new Map()

// request 拦截器
axios.interceptors.request.use(
  function (config) {
    // 登录接口就不拦截了
    if (config.url.includes('login')) {
      return config
    }
    const accessToken = store.get('access_token')
    config.headers = {
      Authorization: `Bearer ${accessToken}`,
    }
    return config
  },
  function (error) {
    return Promise.reject(error)
  }
)

// response 拦截器
axios.interceptors.response.use(
  function (response) {
    // 2xx 范围内的状态码都会触发该函数。
    // Token 无效跳转到登录页面
    if (response.data.code === CODE_INVALID_TOKEN) {
      history.push('/login')
    }
    return response
  },
  function (error) {
    // 超出 2xx 范围的状态码都会触发该函数。
    // 对响应错误做点什么
    return Promise.reject(error)
  }
)

export default function request(options) {
  let { data, url } = options
  const cloneData = cloneDeep(data)

  try {
    let domain = ''
    const urlMatch = url.match(/[a-zA-z]+:\/\/[^/]*/)
    if (urlMatch) {
      ;[domain] = urlMatch
      url = url.slice(domain.length)
    }

    const match = parse(url)
    url = compile(url)(data)

    for (const item of match) {
      if (item instanceof Object && item.name in cloneData) {
        delete cloneData[item.name]
      }
    }
    url = domain + url
  } catch (e) {
    message.error(e.message)
  }

  options.url = url
  options.cancelToken = new CancelToken((cancel) => {
    window.cancelRequest.set(Symbol(Date.now()), {
      pathname: window.location.pathname,
      cancel,
    })
  })

  return axios(options)
    .then((response) => {
      const { statusText, status, data } = response
      let result = {}
      if (typeof data === 'object') {
        result = data
        if (Array.isArray(data)) {
          result.list = data
        }
      } else {
        result.data = data
      }
      console.log(result, '&&&&&')
      return Promise.resolve({
        success: true,
        message: statusText,
        statusCode: status,
        ...result,
      })
    })
    .catch((error) => {
      const { response, message } = error
      console.log(response, message, '.......')
      if (String(message) === CANCEL_REQUEST_MESSAGE) {
        return {
          success: false,
        }
      }

      let msg
      let statusCode

      if (response && response instanceof Object) {
        const { data, statusText } = response
        statusCode = response.status
        msg = data.msg || data.message || statusText
      } else {
        statusCode = 600
        msg = error.message || 'Network Error'
      }

      /* eslint-disable */
      return Promise.reject({
        success: false,
        statusCode,
        message: msg,
      })
    })
}
