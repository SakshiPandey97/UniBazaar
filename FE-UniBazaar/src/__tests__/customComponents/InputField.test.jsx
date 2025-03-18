import { render, screen, fireEvent } from "@testing-library/react";
import InputField from "../../customComponents/InputField";
import { Formik, Form } from "formik";
import { vi } from "vitest";

describe("InputField Component", () => {
  const renderInputField = (props) => {
    render(
      <Formik initialValues={{ testField: "" }} onSubmit={() => {}}>
        <Form>
          <InputField name="testField" data_testid="input-field" {...props} />
        </Form>
      </Formik>
    );
  };

  test("renders the input field correctly", () => {
    renderInputField({ type: "text", placeholder: "Enter text" });

    // Check if the input field is rendered
    const input = screen.getByTestId("input-field");
    expect(input).toBeInTheDocument();

    // Check placeholder text
    expect(input).toHaveAttribute("placeholder", "Enter text");
  });

  test("disables input field when isSubmitting is true", () => {
    renderInputField({ type: "text", isSubmitting: true });

    // Find the input and check if it is disabled
    const input = screen.getByTestId("input-field");
    expect(input).toBeDisabled();
  });

  test("handles onFocus and onBlur events", () => {
    const onFocusMock = vi.fn();
    const onBlurMock = vi.fn();

    renderInputField({
      type: "text",
      additionalProps: { onFocus: onFocusMock, onBlur: onBlurMock },
    });

    const input = screen.getByTestId("input-field");

    // Simulate focus and blur events
    fireEvent.focus(input);
    expect(onFocusMock).toHaveBeenCalledTimes(1);

    fireEvent.blur(input);
    expect(onBlurMock).toHaveBeenCalledTimes(1);
  });

  test("displays error message when there is a Formik validation error", async () => {
    render(
      <Formik
        initialValues={{ testField: "" }}
        initialErrors={{ testField: "Required field" }}
        initialTouched={{ testField: true }}
        onSubmit={() => {}}
      >
        <Form>
          <InputField name="testField" data_testid="input-field" />
        </Form>
      </Formik>
    );

    // Check if the error message is displayed
    expect(screen.getByText("Required field")).toBeInTheDocument();
  });
});
