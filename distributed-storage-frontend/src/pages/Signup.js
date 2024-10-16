import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

const Signup = () => {
  const [Email, setEmail] = useState('');
  const [OTP, setOtp] = useState('');
  const [Password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleSignup = async (e) => {
    e.preventDefault();
    try {
      // After the OTP is sent, ask for it
      const otpResponse = await axios.post('http://localhost:8080/api/signup', { Email, OTP, Password });
      if (otpResponse.status === 200) {
        navigate('/login');
      }
    } catch (error) {
      console.error('Signup failed', error);
      alert('Failed to register');
    }
  };

  return (
    <div>
      <h2>Sign Up</h2>
      <form onSubmit={handleSignup}>
        <div>
          <label>Email:</label>
          <input
            type="email"
            value={Email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <label>OTP:</label>
          <input
            type="text"
            value={OTP}
            onChange={(e) => setOtp(e.target.value)}
            required
          />
        </div>
        <div>
          <label>Password:</label>
          <input
            type="text"
            value={Password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit">Sign Up</button>
      </form>
    </div>
  );
};

export default Signup;
