import React, { useEffect, useState } from 'react';
import axios from 'axios';
import FileUpload from './FileUpload';
import './Dashboard.css';

const Dashboard = () => {
  const [files, setFiles] = useState([]);
  const [errorMessage, setErrorMessage] = useState('');

  const fetchFiles = async () => {
    try {
      document.cookie = `token=${localStorage.getItem('token')}; path=/`;

      const response = await axios.get('http://localhost:8080/api/files', {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
        withCredentials : true,
      });
      setFiles(response.data);
    } catch (error) {
      console.error('Error fetching files', error);
      setErrorMessage('Error fetching files.');
    }
  };

  const deleteFile = async (fileId) => {
    try {
      document.cookie = `token=${localStorage.getItem('token')}; path=/`;

      await axios.delete(`http://localhost:8080/api/delete/${fileId}`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
        withCredentials : true,
      });
      fetchFiles(); // Refresh the file list after deletion
    } catch (error) {
      console.error('Error deleting file', error);
      setErrorMessage('Error deleting file.');
    }
  };

  useEffect(() => {
    fetchFiles();
  }, []);

  return (
    <div className='Dashcontainer'>
      <h2 className='DashHeading2'>My Files</h2>
      {errorMessage && <p style={{ color: 'red' }}>{errorMessage}</p>}
      <FileUpload onUpload={fetchFiles} />
      <ul className = 'Dashul'>
        {files.map((file) => (
          <li class='Dashli' key={file.fid}>
            <a className='Dasha' href={`http://localhost:8080/api/download/${file.email}/${file.filename}`} download>
              {file.filename}
            </a>
            <button className="DelDashbutton" onClick={() => deleteFile(file.fid)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default Dashboard;
