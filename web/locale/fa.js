// eslint-disable-next-line import/no-extraneous-dependencies
const { Organization } = require('@aasaam/information');

module.exports = {
  title: 'محافظ وب‌سایت',
  organization: Organization.fa.description,
  oneMoreStep: 'بررسی درخواست',
  pleaseWait: 'لطفا صبر کنید...',
  submit: 'ارسال',
  noScript:
    'مرورگر شما از جاوا اسکریپت پشتیبانی نمی‌کند. برای ادامه جاوا اسکریپت مرورگر خود را فعال کنید.',
  pleaseEnterCaptcha:
    'برای ادامه کد امنیتی مشاهده شده در تصویر را وارد نمایید. در صورت ناخوانا بودن تصویر می‌توانید با کلیک بر روی <strong>کد جدید</strong> تصویر جدید را امتحان کنید.',
  captcha: 'کد امنیتی',
  pleaseEnterOtp: 'برای ادامه رمز یکبار مصرف خود را وارد نمایید.',
  otp: 'رمز یکبار مصرف',

  pleaseEnterUserPass:
    'برای ادامه نام کاربری،‌ کلمه عبور رو به همراه کد امنیتی موجود در تصویر را وارد نمایید.',
  pleaseEnterYourMobilePhone:
    'برای ادامه شماره تلفن همراه خود را جهت اعتبار سنجی وارد نمایید.',

  invalidUserNameOrPassword: 'نام کاربری و یا کلمه عبور صحیح نمی‌باشد.',
  invalidCaptchaCode: 'کد امنیتی صحیح نمی‌باشد.',
  invalidOTP: 'رمز یکبار مصرف صحیح نمی‌باشد.',
  invalidPhoneNumber: 'شماره تلفن همراه صحیح نمی‌باشد.',

  updateCaptchaImage: 'کد جدید',
  username: 'نام کاربری',
  password: 'رمز عبور',
  customMessageCaptcha: '',
  customMessageJS: '',
  customMessageOTP:
    'شما درخواست استفاده یکی از سرویس های خصوصی را دارید که باید از پیش کد امنیتی رمز یکبار مصرف را داشته باشید.',
  customMessageUsernamePassword: '',
  customMessageSMS: ''
};
