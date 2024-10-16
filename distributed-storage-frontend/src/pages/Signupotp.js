import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

const Signupotp = () => {
  const [Email, setEmail] = useState('');
  const navigate = useNavigate();

  const handleSignup_top = async (e) => {
    e.preventDefault();
    try {
      // Send OTP first to the user's email
      // const otpResponse = await axios.post('http://localhost:8080/api/send-otp', { Email });
      const otpResponse = await axios.post('http://localhost:8080/api/send-otp', {
        email: Email, // Ensure the key matches what the backend expects
      }, {
        headers: {
          'Content-Type': 'application/json', // Explicitly set the content type
        }
      });
      

      if (otpResponse.status === 200) {
        navigate('/signup');
      }
    } catch (error) {
      console.error('Error in sending the OTP', error);
      alert('Failed to send OTP');
    }
  };

  return (
    <div>
      <h2>Request an OTP for Signing Up:</h2>
      <form onSubmit={handleSignup_top}>
        <div>
          <label>Email:</label>
          <input
            type="email"
            value={Email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <button type="submit">Sign Up</button>
      </form>
    </div>
  );
};

export default Signupotp;
