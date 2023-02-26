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

  queryPostList: '/posts',

  queryDemand: '/demand/:id',
  queryDemandList: '/demand',
  createDemand: 'POST /demand',
  updateDemand: 'POST /demand/:id',
  publishDemand: 'POST /demand/:id/publish',
  removeDemand: 'DELETE /demand/:id',
  removeDemandList: 'POST /demands/delete',

  queryTaskList: '/demand/:id/task',

  queryDashboard: '/dashboard',
}
