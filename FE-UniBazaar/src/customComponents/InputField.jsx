import React from "react";
import { Field, ErrorMessage } from "formik";

function InputField({
  data_testid,
  label,
  name,
  type,
  isSubmitting,
  placeholder,
  onFocus,
  onBlur,
  onChange,
  additionalProps = {},
}) {
  return (
    <div className="flex flex-col mb-4 relative">
      <div className="flex justify-center">
        <Field
          data_testid={data_testid}
          placeholder={placeholder}
          type={type}
          id={name}
          name={name}
          disabled={isSubmitting}
          className="mt-1 p-2 w-4/5 border rounded-md focus:outline-none focus:ring focus:ring-[#6D9886] disabled:bg-gray-200 bg-white"
          onFocus={onFocus}
          onBlur={onBlur}
          onChange={onChange}
          {...additionalProps}
        />
      </div>
      <ErrorMessage name={name} component="div" className="text-red-500 text-sm mt-1 ml-[50px]" />
    </div>
  );
}

export default InputField;
