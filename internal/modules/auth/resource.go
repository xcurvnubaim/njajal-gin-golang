package auth

func templateSendEmail(otp string) string {
	return `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta http-equiv="X-UA-Compatible" content="IE=edge">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Your OTP Code</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					margin: 0;
					padding: 0;
					background-color: #f4f4f4;
					color: #333;
				}
				.email-container {
					max-width: 600px;
					margin: 20px auto;
					background-color: #ffffff;
					border-radius: 8px;
					box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
					overflow: hidden;
				}
				.email-header {
					background-color: #4CAF50;
					color: #ffffff;
					padding: 20px;
					text-align: center;
				}
				.email-body {
					padding: 20px;
					line-height: 1.6;
				}
				.otp-code {
					font-size: 24px;
					font-weight: bold;
					color: #4CAF50;
				}
				.email-footer {
					background-color: #f4f4f4;
					text-align: center;
					padding: 10px;
					font-size: 12px;
					color: #666;
				}
				a {
					color: #4CAF50;
					text-decoration: none;
				}
			</style>
		</head>
		<body>
			<div class="email-container">
				<div class="email-header">
					<h1>Your OTP Code</h1>
				</div>
				<div class="email-body">
					<p>Hello,</p>
					<p>Thank you for using our service. Here is your OTP code to proceed:</p>
					<p class="otp-code">` + otp + `</p>
					<p>Please use this code within the next 10 minutes. If you didn't request an OTP, please ignore this email.</p>
				</div>
				<div class="email-footer">
					<p>&copy; 2024 Your Company. All rights reserved.</p>
					<p>Need help? <a href="mailto:support@yourcompany.com">Contact Support</a></p>
				</div>
			</div>
		</body>
		</html>
	`
}