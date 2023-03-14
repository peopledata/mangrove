module.exports = {
  siteName: 'PeopleData Admin',
  copyright: 'PeopleData Admin ©2023',
  logoPath: '/logo.svg',
  apiPrefix: '/admin/api/v1',
  fixedHeader: true, // sticky primary layout header

  /* Layout configuration, specify which layout to use for route. */
  layouts: [
    {
      name: 'primary',
      include: [/.*/],
      exclude: [/(\/(en|zh))*\/login/],
    },
  ],

  /* I18n configuration, `languages` and `defaultLanguage` are required currently. */
  // i18n: {
  //   /* Countrys flags: https://www.flaticon.com/packs/countrys-flags */
  //   languages: [
  //     {
  //       key: 'zh',
  //       title: '中文',
  //       flag: '/china.svg',
  //     },
  //     {
  //       key: 'en',
  //       title: 'English',
  //       flag: '/america.svg',
  //     },
  //   ],
  //   defaultLanguage: 'zh',
  // },
}
