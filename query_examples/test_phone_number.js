var phone = require('node-phonenumber')

var phoneUtil = phone.PhoneNumberUtil.getInstance();
var phoneNumber = phoneUtil.parse('77053237001','MY');
var toNumber = phoneUtil.format(phoneNumber, phone.PhoneNumberFormat.E164);

console.log(toNumber);
