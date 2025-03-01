import React from "react";
import { XCircle,CheckCircle } from "lucide-react";

function PasswordReqBox({password}) {
      const requirements = [
        { label: "At least 6 characters", valid: password.length >= 6 },
        { label: "At least one uppercase letter", valid: /[A-Z]/.test(password) },
        { label: "At least one lowercase letter", valid: /[a-z]/.test(password) },
        { label: "At least one number", valid: /\d/.test(password) },
        { label: "At least one special character (@, $, !, %, *, ?, &)", valid: /[@$!%*?&]/.test(password) },
      ];
    
  return (
    <div className="absolute left-[50px] top-[70px] bg-white shadow-lg border rounded-lg p-3 w-4/5 z-10">
      <p className="text-sm font-medium text-gray-700 mb-1">
        Password must contain:
      </p>
      <ul className="text-sm">
        {requirements.map((req, index) => (
          <li key={index} className="flex items-center gap-2">
            {req.valid ? (
              <CheckCircle className="text-green-500 w-4 h-4" />
            ) : (
              <XCircle className="text-red-500 w-4 h-4" />
            )}
            <span className={req.valid ? "text-green-500" : "text-gray-700"}>
              {req.label}
            </span>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default PasswordReqBox;
