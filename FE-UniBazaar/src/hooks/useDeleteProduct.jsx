import { useState, useCallback } from 'react';
import { deleteProductAPI } from '../api/productAxios';

const useDeleteProduct = () => {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);
    const [isSuccess, setIsSuccess] = useState(false);

    const deleteProduct = useCallback(async (userId, productId) => {
        setIsLoading(true);
        setError(null);
        setIsSuccess(false);

        try {
            await deleteProductAPI(userId, productId);
            setIsSuccess(true);
        } catch (err) {
            setError(err);
            setIsSuccess(false);
        } finally {
            setIsLoading(false);
        }
    }, []);

    return {
        isLoading,
        error,
        isSuccess,
        deleteProduct,
    };
};

export default useDeleteProduct;