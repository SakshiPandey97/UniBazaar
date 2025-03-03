import React from "react";

export default function Spinner() {
  return (
    <div className="spinner-container" data-testid="spinner-container" role="status">
      <div className="spinner"></div>
    </div>
  );
};
