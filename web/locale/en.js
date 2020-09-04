// eslint-disable-next-line import/no-extraneous-dependencies
const { Organization } = require('@aasaam/information');

module.exports = {
  title: 'Website Protection',
  organization: Organization.en.description,
  oneMoreStep: 'One more step',
  submit: 'Submit',
  pleaseWait: 'Please wait...',
  noScript:
    'JavaScript not supported by your browser. Try to enable JavaScript in your browser to continue.',
  pleaseEnterCaptcha:
    'Please enter the image content for continue access website. If the image is not clear click on <strong>New Image</strong> for another try.',
  captcha: 'Security code',
  pleaseEnterOtp:
    'Please enter your OTP(One time password) for continue access website.',
  otp: 'OTP (One Time Password)',

  pleaseEnterUserPass:
    'Please enter your username, password and also type security code of image.',
  pleaseEnterYourMobilePhone:
    'Please enter your mobile phone number for get verification code.',

  invalidUserNameOrPassword: 'Username or password is not correct.',

  invalidCaptchaCode: 'Security code is not correct',
  invalidOTP: 'OTP (One Time Password) is not correct',
  invalidPhoneNumber: 'Mobile phone number is not correct',

  updateCaptchaImage: 'New image',
  username: 'Username',
  password: 'Password',
  customMessageCaptcha: '',
  customMessageJS: '',
  customMessageOTP:
    'You are try to connect to private service, so you must have OTP secret for continue usage of this service.',
  customMessageUsernamePassword: '',
  customMessageSMS: ''
};
