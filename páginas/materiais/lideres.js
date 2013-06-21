$(function() {
	confirmacao()
});

function confirmacao() {
	$('#enviar').click(function(e) {
		var turnos = $.map($("#turnos > input:checked + label"), function(x) {
			return $(x).text()
		});

		$("#confirmacao-turnos").html(turnos.join("<br/>"));
		$("#confirmacao-zona").html($("#zona > :selected").html());

		e.preventDefault();
		$('#confirmacao').dialog({
			resizable: false,
			dialogClass: "dialog",
			height: 330,
			width: 680,
			modal: true,
			title: "Confirmação",
			buttons: {
				"Não": function() {
					$(this).dialog("close");
				},
				"Sim": function() {
					$("#formulario").submit();
					$(this).dialog("close");
				}
			}
		});
	});
}
