// import './sign.css';
import { useState } from 'react';


// react input component for 
// geting username and password
// and sending it to the server
function SignForm() {
  const [username, setUsername] = useState('murat');
  const [password, setPassword] = useState('qwerty');

  async function sendData(data) {
    const response = await fetch("http://localhost:8081/register", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    });
    const body = await response.json();
    return body;
  }

  const handleSubmit = (e) => {
    e.preventDefault();

    const data = {
        username: username,
        password: password
    }

    // clear the input fields
    setUsername('');
    setPassword('');

    // send the data to the server
    sendData(data);
  }

  return (
    <form onSubmit={handleSubmit}>
      <input type="text" placeholder="username" value={username} onChange={(e) => setUsername(e.target.value)} />
      <input type="password" placeholder="password" value={password} onChange={(e) => setPassword(e.target.value)} />
      <input type="submit" value="register" />
    </form>
  );
}


export default SignForm;
