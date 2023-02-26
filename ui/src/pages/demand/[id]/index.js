import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'umi'
import { Page } from 'components'
import styles from './index.less'

@connect(({ demandDetail }) => ({ demandDetail }))
class DemandDetail extends PureComponent {
  render() {
    const { demandDetail } = this.props
    const { data } = demandDetail
    const content = []
    for (let key in data) {
      if ({}.hasOwnProperty.call(data, key)) {
        content.push(
          <div key={key} className={styles.item}>
            <div>{key}</div>
            <div>{String(data[key])}</div>
          </div>
        )
      }
    }
    return (
      <Page inner>
        <div className={styles.content}>{content}</div>
      </Page>
    )
  }
}

DemandDetail.propTypes = {
  demandDetail: PropTypes.object,
}

export default DemandDetail
