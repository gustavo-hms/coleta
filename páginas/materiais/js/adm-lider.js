$(function() {
	confirmacao();
	atualizarEsquinas();
});

function confirmacao() {
	$('#alterar').click(function(e) {
		var idDaZona = $("#zona > :selected").val();
//		obterEsquinas(idDaZona)

		e.preventDefault();
		$('#selecionar-esquina').dialog({
			resizable: false,
			dialogClass: "dialog",
			height: 330,
			width: 480,
			modal: true,
			title: "Alterar esquina",
			buttons: {
				"Cancelar": function() {
					$(this).dialog("close");
				},
				"Ok": function() {
					$("#esquina").val($(".selecao-esquina").val())
					$(this).dialog("close");
				}
			}
		});
	});
}

function obterEsquinas(idDaZona) {
	$.get("../zona/" + idDaZona + "/esquinas", function(esquinas) {
		var lista = []
		for (i in esquinas) {
			id = esquinas[i].id
			element = "<label for='" + id + "'>" +
				"<input type='radio' class='selecao-esquina' name='selecao-esquina' id='" +
					esquinas[i].id + "' value='" + esquinas[i].id + "'>" +
				esquinas[i].cruzamento + "</label>"

			lista.push(element)

			$("#esquinas-exibidas").html(lista)
		}
	})
}

function atualizarEsquinas() {
	$("#selecao-zona").change(function() {
		var idDaZona = $("#zona > :selected").val();
		obterEsquinas(idDaZona)
	})
}
