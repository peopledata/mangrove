export const ROLE_TYPE = {
  ADMIN: 'admin',
  DEFAULT: 'admin',
  DEVELOPER: 'developer',
}

export const CANCEL_REQUEST_MESSAGE = 'cancel request'

// 请求成功状态码
export const CODE_SUCCESS = 1000
// 无效的认证Token
export const CODE_INVALID_TOKEN = 1008

// 草稿：1 已发布：2 已完成：3  运行中：4  已下架：5
export const DEMAND_STATUS = {
  INIT: 1,
  PUBLISHED: 2,
  COMPLETED: 3,
  CLOSED: 4,
  PUBLISHING: 5,
  PUBLISHFAIL: 6,
  RUNNING: 7,
}

export const DEMAND_STATUS_MAP = [
  {
    value: DEMAND_STATUS.INIT,
    color: 'default',
    label: '草稿',
  },
  {
    value: DEMAND_STATUS.PUBLISHED,
    color: 'success',
    label: '已发布',
  },
  {
    value: DEMAND_STATUS.COMPLETED,
    label: '已完成',
    color: 'processing',
  },
  {
    value: DEMAND_STATUS.PUBLISHING,
    label: '发布中',
    color: 'warning',
  },
  {
    value: DEMAND_STATUS.RUNNING,
    label: '运行中',
    color: 'warning',
  },
  {
    value: DEMAND_STATUS.CLOSED,
    label: '已下架',
    color: 'error',
  },
]

// 获取status对应的标签
export const getDemandStatusLabel = (status) => {
  const item = DEMAND_STATUS_MAP.find((item) => item.value === status)
  return item ? item.label : status
}

export const getDemandStatusColor = (status) => {
  const item = DEMAND_STATUS_MAP.find((item) => item.value === status)
  return item ? item.color : 'default'
}

export const DEMAND_CATEGORY = {
  BANK: 'bank',
  CONSUM: 'consum',
  TELECOM: 'telecom',
}

export const DEMAND_APP = {
  WEBANK: 'webank',
  JD: 'jd',
}

export const DEMAND_CATEGORY_MAP = [
  {
    value: DEMAND_CATEGORY.BANK,
    label: '金融',
  },
  {
    value: DEMAND_CATEGORY.CONSUM,
    label: '消费',
  },
  {
    value: DEMAND_CATEGORY.TELECOM,
    label: '电信',
  },
]

export const getDemandCategoryLabel = (category) => {
  const item = DEMAND_CATEGORY_MAP.find((item) => item.value === category)
  return item ? item.label : category
}

// 运行中：1 运行成功：2 运行失败：3
export const TASK_STATUS = {
  RUNNING: 1,
  SUCCESS: 2,
  FAILED: 3,
}

export const TASK_STATUS_MAP = [
  {
    value: TASK_STATUS.RUNNING,
    color: 'warning',
    label: '运行中',
  },
  {
    value: TASK_STATUS.SUCCESS,
    color: 'success',
    label: '运行成功',
  },
  {
    value: TASK_STATUS.FAILED,
    label: '运行失败',
    color: 'error',
  },
]

// 获取status对应的标签
export const getTaskStatusLabel = (status) => {
  const item = TASK_STATUS_MAP.find((item) => item.value === status)
  return item ? item.label : status
}

export const getTaskStatusColor = (status) => {
  const item = TASK_STATUS_MAP.find((item) => item.value === status)
  return item ? item.color : 'default'
}

export const DEMAND_APP_MAP = [
  {
    value: DEMAND_APP.WEBANK,
    label: '微众银行',
  },
  {
    value: DEMAND_APP.JD,
    label: '京东',
  },
]

export const getDemandAppLabel = (category) => {
  const item = DEMAND_APP_MAP.find((item) => item.value === category)
  return item ? item.label : category
}
