import React, { PureComponent } from 'react'
import ReactQuill from 'react-quill'
import 'react-quill/dist/quill.snow.css'
import PropTypes from 'prop-types'
import {
  Divider,
  Row,
  Col,
  Form,
  Input,
  InputNumber,
  Select,
  Modal,
  DatePicker,
} from 'antd'
import { Trans } from '@lingui/macro'
import { t } from '@lingui/macro'
import { DEMAND_CATEGORY_MAP, DEMAND_APP_MAP } from 'utils/constant'

const FormItem = Form.Item

const formItemLayout = {
  labelCol: {
    span: 10,
  },
  wrapperCol: {
    span: 20,
  },
  style: { marginBottom: '.6rem' },
}

class DemandModal extends PureComponent {
  formRef = React.createRef()

  checkPurpose = (_, value) => {
    if (Object.keys(value).length === 0) {
      if (value.purpose) {
        return Promise.resolve()
      }
      return Promise.reject(new Error('数据用途内容不能为空！'))
    }
    return Promise.resolve()
  }

  checkAgreement = (_, value) => {
    if (Object.keys(value).length === 0) {
      if (value.agreement) {
        return Promise.resolve()
      }
      return Promise.reject(new Error('数据协议内容不能为空！'))
    }
    return Promise.resolve()
  }

  handleOk = () => {
    const { item = {}, onOk } = this.props

    this.formRef.current
      .validateFields()
      .then((values) => {
        const data = {
          ...values,
          demand_id: item.demand_id,
        }
        onOk(data)
      })
      .catch((err) => {
        console.log(err)
      })
  }

  render() {
    const { item = {}, onOk, form, ...modalProps } = this.props
    return (
      <Modal {...modalProps} onOk={this.handleOk}>
        <Form
          ref={this.formRef}
          name="control-ref"
          initialValues={{
            ...item,
          }}
          layout="vertical"
        >
          <Divider></Divider>
          <Row>
            <Col>
              <h3>{t`Basic Info`}</h3>
            </Col>
          </Row>
          <Row gutter={24}>
            <Col span={12}>
              <FormItem
                name="name"
                rules={[{ required: true }]}
                label={t`Demand Name`}
                hasFeedback
                {...formItemLayout}
              >
                <Input />
              </FormItem>
            </Col>
            <Col span={12}>
              <FormItem
                name="valid_at"
                rules={[{ required: true }]}
                label={t`Demand Valid Time`}
                hasFeedback
                {...formItemLayout}
              >
                <DatePicker showTime />
              </FormItem>
            </Col>
          </Row>
          <Row gutter={24}>
            <Col span={12}>
              <FormItem
                name="category"
                rules={[{ required: true }]}
                label={t`Data Category`}
                hasFeedback
                {...formItemLayout}
              >
                <Select
                  style={{
                    width: 120,
                  }}
                  placeholder={t`Select Category`}
                  options={DEMAND_CATEGORY_MAP}
                />
              </FormItem>
            </Col>
            <Col span={12}>
              <FormItem
                name="content"
                rules={[{ required: true }]}
                label={t`Data Content`}
                hasFeedback
                {...formItemLayout}
              >
                <Select
                  style={{
                    width: 120,
                  }}
                  placeholder={t`Select APP`}
                  options={DEMAND_APP_MAP}
                />
              </FormItem>
            </Col>
          </Row>
          <Row gutter={24}>
            <Col span={22}>
              <FormItem
                name="brief"
                rules={[{ required: true }]}
                label={t`Brief`}
                hasFeedback
              >
                <Input.TextArea />
              </FormItem>
            </Col>
          </Row>
          <Row gutter={24}>
            <Col span={12}>
              <FormItem
                name="need_users"
                rules={[{ required: true }]}
                label={t`Needs Users`}
                hasFeedback
                {...formItemLayout}
              >
                <InputNumber min={1} />
              </FormItem>
            </Col>
            <Col span={12}>
              <FormItem
                name="use_times"
                rules={[{ required: true }]}
                label={t`Data Used Times`}
                hasFeedback
                {...formItemLayout}
              >
                <InputNumber min={1} />
              </FormItem>
            </Col>
          </Row>

          <Row>
            <Col span={22}>
              <FormItem
                name="purpose"
                label={t`Data Usage`}
                rules={[
                  {
                    validator: this.checkPurpose,
                  },
                ]}
                style={{ marginBottom: '.6rem' }}
              >
                <ReactQuill style={{ height: '100px' }} theme="snow" />
              </FormItem>
            </Col>
          </Row>

          <Row style={{ marginTop: '3rem' }}>
            <Col>
              <h3>{t`Algorithms and protocols`}</h3>
            </Col>
          </Row>

          <Row gutter={24} style={{ marginBottom: '3rem' }}>
            <Col span={12}>
              <FormItem
                name="algorithm"
                rules={[{ required: true }]}
                label={t`Algorithm file image address`}
                hasFeedback
                {...formItemLayout}
              >
                <Input />
              </FormItem>
            </Col>
            <Col span={12}>
              <FormItem
                name="command"
                label={t`Algorithm Command`}
                hasFeedback
                {...formItemLayout}
              >
                <Input />
              </FormItem>
            </Col>
          </Row>

          <Row style={{ marginBottom: '3rem' }}>
            <Col span={22}>
              <FormItem
                name="agreement"
                label={t`Agreement Content`}
                rules={[
                  {
                    validator: this.checkAgreement,
                  },
                ]}
                style={{ marginBottom: '.6rem' }}
              >
                <ReactQuill style={{ height: '100px' }} theme="snow" />
              </FormItem>
            </Col>
          </Row>
        </Form>
      </Modal>
    )
  }
}

DemandModal.propTypes = {
  type: PropTypes.string,
  item: PropTypes.object,
  onOk: PropTypes.func,
}

export default DemandModal
