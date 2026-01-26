import java.math.BigDecimal;
import java.math.RoundingMode;
import java.text.NumberFormat;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Currency;

interface Payable {
    boolean processPayment();
}

public class PaymentTransaction implements Payable {

    private String transactionId;
    private BigDecimal amount;
    private String currency;
    private String merchantId;
    private Date timestamp;

    public PaymentTransaction(String transactionId, BigDecimal amount,
                              String currency, String merchantId, Date timestamp) {
        this.transactionId = transactionId;
        this.amount = amount;
        this.currency = currency;
        this.merchantId = merchantId;
        this.timestamp = timestamp;
    }

    public String getTransactionId() {
        return transactionId;
    }

    public void setTransactionId(String transactionId) {
        this.transactionId = transactionId;
    }

    public BigDecimal getAmount() {
        return amount;
    }

    public void setAmount(BigDecimal amount) {
        this.amount = amount;
    }

    public String getCurrency() {
        return currency;
    }

    public void setCurrency(String currency) {
        this.currency = currency;
    }

    public String getMerchantId() {
        return merchantId;
    }

    public void setMerchantId(String merchantId) {
        this.merchantId = merchantId;
    }

    public Date getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(Date timestamp) {
        this.timestamp = timestamp;
    }

    public BigDecimal calculateFee(double feePercentage) {
        if (amount == null) {
            return BigDecimal.ZERO;
        }
        BigDecimal percentage = BigDecimal.valueOf(feePercentage / 100);
        return amount.multiply(percentage).setScale(2, RoundingMode.HALF_UP);
    }

    public String getDisplayInfo() {
        SimpleDateFormat dateFormat = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        NumberFormat currencyFormatter = NumberFormat.getCurrencyInstance();

        try {
            currencyFormatter.setCurrency(Currency.getInstance(currency));
        } catch (IllegalArgumentException e) {
            // Fallback if currency code is invalid
        }

        return String.format(
                "Transaction ID: %s%n" +
                        "Amount: %s%n" +
                        "Currency: %s%n" +
                        "Merchant ID: %s%n" +
                        "Timestamp: %s",
                transactionId,
                amount != null ? currencyFormatter.format(amount) : "N/A",
                currency,
                merchantId,
                timestamp != null ? dateFormat.format(timestamp) : "N/A"
        );
    }

    // Implement Payable interface
    @Override
    public boolean processPayment() {
        // Implementation logic for processing payment
        if (amount == null || amount.compareTo(BigDecimal.ZERO) <= 0) {
            return false;
        }
        if (transactionId == null || transactionId.isEmpty()) {
            return false;
        }

        // Simulate payment processing
        System.out.println("Processing payment: " + transactionId);
        // In real implementation, this would connect to payment gateway

        return true;
    }

    // Optional: Override equals, hashCode, and toString
    @Override
    public String toString() {
        return "PaymentTransaction{" +
                "transactionId='" + transactionId + '\'' +
                ", amount=" + amount +
                ", currency='" + currency + '\'' +
                ", merchantId='" + merchantId + '\'' +
                ", timestamp=" + timestamp +
                '}';
    }
}