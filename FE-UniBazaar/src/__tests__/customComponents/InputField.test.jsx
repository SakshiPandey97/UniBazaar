import { render, screen, fireEvent } from "@testing-library/react";
import InputField from "../../customComponents/InputField";
import { Formik, Form } from "formik";
import { vi } from "vitest";

describe("InputField Component", () => {
  const renderInputField = (props) => {
    render(
      <Formik
        initialValues={{ testField: "" }}
        onSubmit={() => {}}
      >
        <Form>
          <InputField name="testField" {...props} />
        </Form>
      </Formik>
    );
  };

  test("renders the input field with label", () => {
    renderInputField({ label: "Test Label", type: "text" });
    
    // Check if the label is rendered correctly
    expect(screen.getByLabelText("Test Label")).toBeInTheDocument();
  });

  test("disables input field when isSubmitting is true", () => {
    renderInputField({ label: "Test Label", type: "text", isSubmitting: true });
    
    // Find the input and check if it is disabled
    const input = screen.getByLabelText("Test Label");
    expect(input).toBeDisabled();
  });

  test("handles onFocus and onBlur events", () => {
    const onFocusMock = vi.fn();
    const onBlurMock = vi.fn();

    renderInputField({ label: "Test Label", type: "text", onFocus: onFocusMock, onBlur: onBlurMock });

    const input = screen.getByLabelText("Test Label");
    
    // Simulate focus and blur events
    fireEvent.focus(input);
    expect(onFocusMock).toHaveBeenCalledTimes(1);

    fireEvent.blur(input);
    expect(onBlurMock).toHaveBeenCalledTimes(1);
  });

  // Uncomment and adjust error handling test if necessary
  // test("displays error message when there is a formik error", () => {
  //   renderInputField({ label: "Test Label", type: "text" });
    
  //   // Simulate form error
  //   const errorMessage = screen.getByText(/testfield/i);
  //   expect(errorMessage).toBeInTheDocument();
  // });

  test("renders password reset link if isPassword is true", () => {
    renderInputField({ label: "Password", type: "password", isPassword: true });
    
    // Check if the forgot password link is rendered
    expect(screen.getByText("Forgot password?")).toBeInTheDocument();
  });

  test("does not render password reset link if isPassword is false", () => {
    renderInputField({ label: "Password", type: "password", isPassword: false });
    
    // Check if the forgot password link is not rendered
    expect(screen.queryByText("Forgot password?")).toBeNull();
  });
});
