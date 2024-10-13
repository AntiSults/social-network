import React from 'react';

const ErrorMessage = (message: string) => {
  return (
    <div>
      <p style={{ color: "red" }}>{message}</p>
    </div>
  );
};

export default ErrorMessage;
