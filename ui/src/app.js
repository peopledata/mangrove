import { message } from 'antd'

// 全局统一处理
export const dva = {
  config: {
    onError(e) {
      e.preventDefault()
      if (e.msg || e.message) {
        message.error(e.msg || e.message)
      } else {
        /* eslint-disable */
        console.error(e)
      }
    },
  },
}
