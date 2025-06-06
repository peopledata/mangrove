import { FrownOutlined } from '@ant-design/icons'
// @ts-ignore
import { Page } from 'components'
// @ts-ignore
import styles from './404.less'

const Error = () => (
  <Page inner>
    <div className={styles.error}>
      <FrownOutlined />
      <h1>404 Not Found</h1>
    </div>
  </Page>
)

export default Error
