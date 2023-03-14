import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'umi'
import { Row, Col, Card } from 'antd'
import { Page } from 'components'
import { NumberCard, Sales } from './components'
import styles from './index.less'

@connect(({ dashboard, loading }) => ({
  dashboard,
  loading,
}))
class Dashboard extends PureComponent {
  render() {
    const { dashboard } = this.props
    const { sales, numbers } = dashboard

    let numberCards = <></>
    if (numbers) {
      numberCards = numbers.map((item, key) => (
        <Col key={key} lg={6} md={12}>
          <NumberCard {...item} />
        </Col>
      ))
    }

    return (
      <Page className={styles.dashboard}>
        <Row gutter={24}>
          {numberCards}
          <Col lg={24} md={24}>
            <Card
              bordered={false}
              bodyStyle={{
                padding: '24px 36px 24px 0',
              }}
            >
              <Sales data={sales} />
            </Card>
          </Col>
        </Row>
      </Page>
    )
  }
}

Dashboard.propTypes = {
  dashboard: PropTypes.object,
  loading: PropTypes.object,
}

export default Dashboard
