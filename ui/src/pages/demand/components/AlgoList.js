import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { Table, Modal, Tag, Drawer } from 'antd'
import { DropOption } from 'components'
import dayjs from 'dayjs'
import { t } from '@lingui/macro'
import { Trans } from '@lingui/macro'
import { Link } from 'umi'
import {
  getDemandStatusLabel,
  getDemandStatusColor,
  getDemandCategoryLabel,
} from 'utils/constant'
import styles from './AlgoList.less'
import { DEMAND_STATUS } from 'utils/constant'

const { confirm } = Modal

class AlgoList extends PureComponent {
  render() {
    const { onOk, dataSource, ...algoListProps } = this.props
    const columns = [
      {
        title: <Trans>Index</Trans>,
        dataIndex: 'index',
        key: 'index',
        width: '10%',
        fixed: 'left',
      },
      {
        title: 'Name',
        dataIndex: 'name',
        key: 'name',
      },
      {
        title: 'Balance',
        dataIndex: 'balance',
        key: 'balance',
      },
    ]
    return (
      <Drawer style={{ zIndex: 1250 }} {...algoListProps}>
        <Table columns={columns} simple dataSource={dataSource} />
      </Drawer>
    )
  }
}

AlgoList.propTypes = {
  location: PropTypes.object,
}

export default AlgoList
