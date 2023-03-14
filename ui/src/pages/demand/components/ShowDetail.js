import React, { PureComponent } from 'react'
import { Drawer, Row, Col, Tabs, Table, Tag } from 'antd'
import { t, Trans } from '@lingui/macro'
import PropTypes from 'prop-types'
import styles from './Detail.less'
import dayjs from 'dayjs'
import { getDemandCategoryLabel } from 'utils/constant'

class ShowDetail extends PureComponent {
  render() {
    const { content, ...drawerProps } = this.props

    return (
      <Drawer {...drawerProps}>
        <div dangerouslySetInnerHTML={{ __html: content }}></div>
      </Drawer>
    )
  }
}

ShowDetail.propTypes = {
  demand: PropTypes.object,
}

export default ShowDetail
