function submitbtn() {
  let nama = document.getElementById('nama').value;
  let email = document.getElementById('email').value;
  let phone = document.getElementById('phone').value;
  let subject = document.getElementById('subject').value;
  let message = document.getElementById('message').value;

  if (nama == '') {
    return alert('Name input fields must be not empty');
  } else if (email == '') {
    return alert('Email input fields must be not empty');
  } else if (phone == '') {
    return alert('Phone input fields must be not empty');
  } else if (subject == '') {
    return alert('Subject input fields must be not empty');
  } else if (message == '') {
    return alert('Message input fields must be not empty');
  }

  const emailReciver = 'dionovalino@gmail.com';

  const a = document.createElement('a');

  a.href = `mailto:${emailReciver}?subject=${subject}&body=Hello my name ${nama}, ${subject}, ${message}`;
  a.click();

  let dataObject = {
    nama: nama,
    email: email,
    phoneNumber: phone,
    subject: subject,
    message: message,
  };
  console.log(dataObject);
  
}

