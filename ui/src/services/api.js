export default {
  queryRouteList: '/routes',

  queryUserInfo: '/user',
  logoutUser: '/logout',
  loginUser: 'POST /login',

  queryUser: '/user/:id',
  queryUserList: '/users',
  updateUser: 'Patch /user/:id',
  createUser: 'POST /user',
  removeUser: 'DELETE /user/:id',
  removeUserList: 'POST /users/delete',

  queryDemandInfo: 'GET /demand/:id/info',
  queryDemandDetail: 'GET /demand/:id/detail',
  queryContractRecords: 'GET /demand/:id/contract_record',
  queryDemandList: '/demand',
  createDemand: 'POST /demand',
  updateDemand: 'POST /demand/:id',
  publishDemand: 'POST /demand/:id/publish',
  removeDemand: 'DELETE /demand/:id',
  removeDemandList: 'POST /demands/delete',

  queryTaskList: '/demand/:id/task',

  queryDashboard: '/dashboard',
}
