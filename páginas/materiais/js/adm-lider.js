$(function() {
	confirmacao();
	atualizarEsquinas();
});

function confirmacao() {
	$('#alterar').click(function(e) {
//		var idDaZona = $("#zona > :selected").val();
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
					$("#esquina").val($(".selecao-esquina:checked").val())
					$("#zona-oculta").val($("#zona").val())
					$("#exibicao-esquina").html($(".selecao-esquina:checked").
						parent().text())
					$(this).dialog("close");
				}
			}
		});
	});
}

function obterEsquinas(idDaZona) {
	$.get("../zona/" + idDaZona + "/esquinas", function(esquinas) {
		var lista = []
		var idAtual = $("#esquina").val()
		for (i in esquinas) {
			id = esquinas[i].id
			var checked = id == idAtual? "checked" : "";
			element = "<label for='" + id + "'>" +
				"<input type='radio' class='selecao-esquina' name='selecao-esquina' id='" +
					esquinas[i].id + "' value='" + esquinas[i].id + "' " + checked + "/> " +
				esquinas[i].cruzamento + "</label><br/>"

			lista.push(element)
		}

		$("#esquinas-exibidas").html(lista)
	})
}

function atualizarEsquinas() {
	$("#zona").change(function() {
		var idDaZona = $("#zona > :selected").val();
		obterEsquinas(idDaZona)
	})
}
