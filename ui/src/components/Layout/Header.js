import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { Layout, Space, Avatar, Dropdown, Menu } from 'antd'
import {
  LogoutOutlined,
  DownOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
} from '@ant-design/icons'
import { Trans } from '@lingui/macro'
import { getLocale, setLocale } from 'utils'
import classnames from 'classnames'
import config from 'config'
import styles from './Header.less'

const { SubMenu } = Menu

class Header extends PureComponent {
  render() {
    const { fixed, avatar, username, collapsed, onCollapseChange, onSignOut } =
      this.props

    const items = [
      {
        label: <Trans>Sign out</Trans>,
        key: 'logout',
        icon: <LogoutOutlined />,
        onClick: function (e) {
          e.key === 'logout' && onSignOut()
        },
      },
    ]

    const rightContent = [
      <Dropdown
        menu={{
          items,
        }}
        trigger={['click']}
      >
        <a className="ant-dropdown-link" onClick={(e) => e.preventDefault()}>
          <Avatar style={{ marginRight: '5px' }} src={avatar} />
          <Space>
            {username}
            <DownOutlined />
          </Space>
        </a>
      </Dropdown>,
    ]

    if (config.i18n) {
      const { languages } = config.i18n
      const language = getLocale()
      const currentLanguage = languages.find((item) => item.key === language)

      rightContent.unshift(
        <Menu
          key="language"
          selectedKeys={[currentLanguage.key]}
          onClick={(data) => {
            setLocale(data.key)
          }}
          mode="horizontal"
        >
          <SubMenu title={<Avatar size="small" src={currentLanguage.flag} />}>
            {languages.map((item) => (
              <Menu.Item key={item.key}>
                <Avatar
                  size="small"
                  style={{ marginRight: 8 }}
                  src={item.flag}
                />
                {item.title}
              </Menu.Item>
            ))}
          </SubMenu>
        </Menu>
      )
    }

    return (
      <Layout.Header
        className={classnames(styles.header, {
          [styles.fixed]: fixed,
          [styles.collapsed]: collapsed,
        })}
        style={{ height: 72, backgroundColor: 'white', paddingInline: 0 }}
        id="layoutHeader"
      >
        <div
          className={styles.button}
          onClick={onCollapseChange.bind(this, !collapsed)}
        >
          {collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
        </div>
        <div className={styles.rightContainer}>{rightContent}</div>
      </Layout.Header>
    )
  }
}

Header.propTypes = {
  fixed: PropTypes.bool,
  user: PropTypes.object,
  menus: PropTypes.array,
  collapsed: PropTypes.bool,
  onSignOut: PropTypes.func,
  notifications: PropTypes.array,
  onCollapseChange: PropTypes.func,
  onAllNotificationsRead: PropTypes.func,
}

export default Header
