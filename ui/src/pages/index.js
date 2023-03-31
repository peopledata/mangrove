import React, { PureComponent } from 'react'
import { Redirect } from 'umi'
import { t } from '@lingui/macro'

class Index extends PureComponent {
  render() {
    return <Redirect to={t`/demand`} />
  }
}

export default Index
