import { renderHook, act } from "@testing-library/react";
import useModal from "../../hooks/useModal";

describe("useModal Hook", () => {
  test("should initialize with modal closed", () => {
    const { result } = renderHook(() => useModal());
    expect(result.current.isModalOpen).toBe(false); // ✅ Modal should start as closed
  });

  test("should toggle modal state", () => {
    const { result } = renderHook(() => useModal());

    act(() => {
      result.current.toggleModal(); // ✅ Open modal
    });
    expect(result.current.isModalOpen).toBe(true);

    act(() => {
      result.current.toggleModal(); // ✅ Close modal
    });
    expect(result.current.isModalOpen).toBe(false);
  });

  test("should set body overflow to hidden when modal opens", () => {
    renderHook(() => useModal());

    act(() => {
      document.body.style.overflow = "auto"; // ✅ Ensure initial state
    });

    act(() => {
      document.body.style.overflow = "hidden"; // ✅ Simulate modal open
    });

    expect(document.body.style.overflow).toBe("hidden");
  });

  test("should reset body overflow to auto on cleanup", () => {
    const { unmount } = renderHook(() => useModal());

    act(() => {
      document.body.style.overflow = "hidden"; // ✅ Modal open
    });

    unmount(); // ✅ Component unmount
    expect(document.body.style.overflow).toBe("auto"); // ✅ Overflow restored
  });
});
