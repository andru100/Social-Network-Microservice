import React from 'react';
import { useEffect } from "react";
import { useNavigate, Link } from 'react-router-dom';
import ChkAuth from './chkAuth';
import Google from '../images/icons/icon-google.png'
import SendData from './SendData';
import UpdatePassword from './UpdatePassword.jsx';
import UpdateUsername from './UpdateUsername';
import UpdateEmail from './UpdateEmail';
import UpdateMobile from './UpdateMobile';
import RequestOTP from './RequestOTP';
import UpdateHybrid from './UpdateHybrid';
import { useAlert } from "react-alert";

export default function RenderUpdateDetails (props) {

	const Navigate = useNavigate();
	const [authtype, setAuthType] = React.useState("init");
	const [username, setUsername] = React.useState(props.username)
	// const [securitylevel, setSecurityLevel] = React.useState(props.securitylevel)
	const [address, setAddress] = React.useState("unassigned")
	const [userDetails, setUserDetails] = React.useState({username: "", password: "", email: "", mobile: ""})
	const [name, setName] = React.useState("John Doe");
	const [email, setEmail] = React.useState("johndoe@gmail.com");
	const [phone, setPhone] = React.useState("555-555-5555");
	const [pagetype, setPageType] = React.useState("init");
  


	console.log("resert request username:", username)

	async function GetOTP (requestType, pagetype) {
		if (requestType == "none") {
			setAuthType(requestType)
			setPageType(pagetype)
			return
		}

		RequestOTP(username, requestType).then((response) => {
			if (( "errors" in response )) {
				console.log("error bak from otp request", response.errors[0].message)
				alertError(response.errors[0].message)
			} else {
				setAddress([response.data.SecureUpdate.MobClue, response.data.SecureUpdate.EmailClue])
				setAuthType(response.data.SecureUpdate.AuthType)
				setPageType(pagetype)
				

			}
		})
	}

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

	function alertError(error){
		const delimiter = '= '
		const start = 2
		const tokens = error.split(delimiter).slice(start)
		const result = tokens.join(delimiter); // those.that
		alert.show(result)
	}

	function InitPage () {
		return (
			<>
				 <div>
					<h1>Account Details</h1>
					<div>
						<label>Username:</label>
						<input type="text" value={name} onChange={handleUpdateName} />
						<button className="login100-form-btn" type="button" onClick={() => GetOTP("unused", "Username")}>
							update username
						</button>
					</div>
					<div>
						<label>Email:</label>
						<input type="text" value={email} onChange={handleUpdateEmail} />
						<button className="login100-form-btn" type="button" onClick={() => GetOTP("unused", "Email")}>
							update email
						</button>
					</div>
					<div>
						<label>Phone:</label>
						<input type="text" value={phone} onChange={handleUpdatePhone} />
						<button className="login100-form-btn" type="button" onClick={() => GetOTP("unused", "Mobile")}>
							update mobile
						</button>
					</div>
					<div>
						<label>Password:</label>
						<input type="text" value={address} onChange={handleUpdateAddress} />
						<button className="login100-form-btn" type="button" onClick={() => GetOTP("unused", "Password")}>
							update password
						</button>
					</div>
				</div>
			</>
		)
	}


	function LandingPage ()	{
		switch (pagetype) { 
			case "init": 
				console.log("found init page")
				return <InitPage/>
				break
			default: return <UpdateHybrid username={username}  address ={address} updatetype = {pagetype} authtype ={authtype} /* setSessionUser={props.setSessionUser()} *//>
		}
	}
	
  return (  
    <>
		<LandingPage/>
	</>


  )
};