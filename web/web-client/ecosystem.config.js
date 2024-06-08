module.exports = {
  apps: [
    {
      name: 'alnitak_web',
      port: '9010', //监听端口
      exec_mode: 'cluster',
      instances: 'max',
      script: './.output/server/index.mjs'
    }
  ]
}
