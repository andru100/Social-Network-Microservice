import React, { useState } from 'react';

export default function AccountDetails() {
  const [name, setName] = useState("John Doe");
  const [email, setEmail] = useState("johndoe@gmail.com");
  const [phone, setPhone] = useState("555-555-5555");
  const [address, setAddress] = useState("123 Main St");

  const handleUpdateName = (event) => {
    setName(event.target.value);
  };

  const handleUpdateEmail = (event) => {
    setEmail(event.target.value);
  };

  const handleUpdatePhone = (event) => {
    setPhone(event.target.value);
  };

  const handleUpdateAddress = (event) => {
    setAddress(event.target.value);
  };

  return (
    <div>
      <h1>Account Details</h1>
      <div>
        <label>Name:</label>
        <input type="text" value={name} onChange={handleUpdateName} />
        <button>Update</button>
      </div>
      <div>
        <label>Email:</label>
        <input type="text" value={email} onChange={handleUpdateEmail} />
        <button>Update</button>
      </div>
      <div>
        <label>Phone:</label>
        <input type="text" value={phone} onChange={handleUpdatePhone} />
        <button>Update</button>
      </div>
      <div>
        <label>Address:</label>
        <input type="text" value={address} onChange={handleUpdateAddress} />
        <button>Update</button>
      </div>
    </div>
  );
}
