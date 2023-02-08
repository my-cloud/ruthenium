$(function () {
    $.ajax({
        url: '/wallet',
        type: 'POST',
        success: function (response) {
            $('#public_key').val(response['public_key']);
            $('#private_key').val(response['private_key']);
            $('#sender_address').val(response['address']);
            console.info(response);
        },
        error: function (error) {
            console.error(error);
        }
    })

    $('#send_money_button').click(function () {
        if (!confirm('Are you sure to send?')) {
            alert('Canceled')
            return
        }

        let transaction_data = {
            'sender_private_key': $('#private_key').val(),
            'sender_address': $('#sender_address').val(),
            'recipient_address': $('#recipient_address').val(),
            'value': $('#send_amount').val(),
        };

        $.ajax({
            url: '/transaction',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify(transaction_data),
            success: function (response) {
                console.info(response);
                if (response === "success") {
                    alert('Send success');
                } else {
                    alert('Send failed: ' + response)
                }
            },
            error: function (response) {
                console.error(response);
                alert('Send failed');
            }
        })
    })

    function refresh_amount() {
        let data = {'address': $('#sender_address').val()}
        $.ajax({
            url: '/wallet/amount',
            type: 'GET',
            data: data,
            success: function (response) {
                $('#wallet_amount').text(response);
            },
            error: function (error) {
                console.error(error)
            }
        })
    }

    function refresh_transactions() {
        $.ajax({
            url: '/transactions',
            type: 'GET',
            success: function (response) {
                $('#transactions_pool').text(JSON.stringify(response, undefined, 4));
            },
            error: function (error) {
                console.error(error)
            }
        })
    }

    function start_validation() {
        $.ajax({
            url: '/validation/start',
            type: 'POST',
            error: function (error) {
                console.error(error)
            }
        })
    }

    function stop_validation() {
        $.ajax({
            url: '/validation/stop',
            type: 'POST',
            error: function (error) {
                console.error(error)
            }
        })
    }

    $('#start_validation').click(function () {
        start_validation();
    });

    $('#stop_validation').click(function () {
        stop_validation();
    });

    setInterval(refresh_amount, 1000)
})