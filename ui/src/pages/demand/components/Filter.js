import React, { Component } from 'react'
import PropTypes from 'prop-types'
import dayjs from 'dayjs'
import { Trans, t } from '@lingui/macro'
import { Button, Row, Col, DatePicker, Form, Input, Select } from 'antd'
import { RedoOutlined } from '@ant-design/icons'
import { DEMAND_STATUS_MAP, DEMAND_CATEGORY_MAP } from 'utils/constant'
import styles from './Filter.less'

const { Search } = Input
const { RangePicker } = DatePicker

const ColProps = {
  xs: 24,
  sm: 12,
  style: {
    marginBottom: 16,
  },
}

const TwoColProps = {
  ...ColProps,
  xl: 96,
}

class Filter extends Component {
  formRef = React.createRef()

  handleFields = (fields) => {
    const { createTime } = fields
    if (createTime && createTime.length) {
      fields.createTime = [
        dayjs(createTime[0]).format('YYYY-MM-DD'),
        dayjs(createTime[1]).format('YYYY-MM-DD'),
      ]
    }
    return fields
  }

  handleSubmit = () => {
    const { onFilterChange } = this.props
    const values = this.formRef.current.getFieldsValue()
    const fields = this.handleFields(values)
    onFilterChange(fields)
  }

  handleReset = () => {
    const fields = this.formRef.current.getFieldsValue()
    for (let item in fields) {
      if ({}.hasOwnProperty.call(fields, item)) {
        if (fields[item] instanceof Array) {
          fields[item] = []
        } else {
          fields[item] = undefined
        }
      }
    }
    this.formRef.current.setFieldsValue(fields)
    this.handleSubmit()
  }

  handleChange = (key, values) => {
    const { onFilterChange } = this.props
    let fields = this.formRef.current.getFieldsValue()
    fields[key] = values
    fields = this.handleFields(fields)
    onFilterChange(fields)
  }

  render() {
    const { onAdd, onRefresh, filter } = this.props
    const { name, status, category } = filter

    let initialCreateTime = []
    if (filter.createTime && filter.createTime[0]) {
      initialCreateTime[0] = dayjs(filter.createTime[0])
    }
    if (filter.createTime && filter.createTime[1]) {
      initialCreateTime[1] = dayjs(filter.createTime[1])
    }

    return (
      <Form
        ref={this.formRef}
        name="control-ref"
        initialValues={{
          name,
          status,
          category,
          createTime: initialCreateTime,
        }}
      >
        <Row gutter={24}>
          <Col {...ColProps} xl={{ span: 3 }} md={{ span: 12 }}>
            <Form.Item name="name">
              <Search
                placeholder={t`Search Demand Name`}
                onSearch={this.handleSubmit}
              />
            </Form.Item>
          </Col>
          <Col
            {...ColProps}
            xl={{ span: 3 }}
            md={{ span: 6 }}
            id="statusSelect"
          >
            <Form.Item name="status">
              <Select
                style={{ width: '100%' }}
                options={DEMAND_STATUS_MAP}
                placeholder={t`Please pick status`}
              />
            </Form.Item>
          </Col>
          <Col
            {...ColProps}
            xl={{ span: 3 }}
            md={{ span: 6 }}
            id="categorySelect"
          >
            <Form.Item name="category">
              <Select
                style={{ width: '100%' }}
                options={DEMAND_CATEGORY_MAP}
                placeholder={t`Please pick category`}
              />
            </Form.Item>
          </Col>
          <Col
            {...TwoColProps}
            xl={{ span: 15 }}
            md={{ span: 24 }}
            sm={{ span: 24 }}
          >
            <Row type="flex" align="middle" justify="space-between">
              <div>
                <Button
                  type="primary"
                  htmlType="submit"
                  className="margin-right"
                  onClick={this.handleSubmit}
                >
                  <Trans>Search</Trans>
                </Button>
                <Button onClick={this.handleReset}>
                  <Trans>Reset</Trans>
                </Button>
              </div>
              <div className={styles.actions}>
                <Button
                  type="primary"
                  ghost
                  icon={<RedoOutlined />}
                  onClick={onRefresh}
                >
                  <Trans>Refresh</Trans>
                </Button>

                <Button danger ghost onClick={onAdd}>
                  <Trans>Create</Trans>
                </Button>
              </div>
            </Row>
          </Col>
        </Row>
      </Form>
    )
  }
}

Filter.propTypes = {
  onAdd: PropTypes.func,
  onRefresh: PropTypes.func,
  filter: PropTypes.object,
  onFilterChange: PropTypes.func,
}

export default Filter
