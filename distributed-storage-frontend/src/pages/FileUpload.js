import React, { useState } from 'react';
import axios from 'axios';
import './Dashboard.css';

const FileUpload = ({ onUpload }) => {
  const [file, setFile] = useState(null);
  const [errormessage, setErrorMessage] = useState('');

  const handleFileUpload = async (e) => {
    e.preventDefault();

    if (!file) {
      setErrorMessage('Please select a file to upload.');
      return;
    }

    const formData = new FormData();
    formData.append('file', file);

    try {
      document.cookie = `token=${localStorage.getItem('token')}; path=/`;
      // console.log(localStorage.getItem('token'));
      await axios.post('http://localhost:8080/api/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',  
        },
        withCredentials : true,
      });
      setErrorMessage(''); // Clear any previous error messages
      onUpload();  // Call the onUpload function to refresh the file list or perform other actions after upload
    } catch (error) {
      if (error.response && error.response.status === 409){
        setErrorMessage('File already exists. Please rename the file or delete the existing one.');
      } else {
        setErrorMessage('File upload failed. Please try again.');
        console.error('File upload failed', error);
      }
    }
  };

  return (
    <form onSubmit={handleFileUpload}>
      <input
        type="file"
        onChange={(e) => setFile(e.target.files[0])}
      />
      <button className='UploadDash' type="submit">Upload File</button>
      {errormessage && <p style = {{color: 'red'}}> {errormessage}</p>}
    </form>
  );
};

export default FileUpload;
