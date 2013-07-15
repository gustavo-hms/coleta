$(function() {
	dialogo();
	atualizarEsquinas();
	registrarFiltro();
});

function dialogo() {
	$('#alterar').click(function(e) {
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
					$("#zona-selecionada").text($("#zona option").filter(":selected").text())
					$("#exibicao-esquina").html($(".selecao-esquina:checked").
						parent().text())
					$(this).dialog("close");
				}
			}
		});
	});
}

var resultados = []
function obterEsquinas(idDaZona) {
	$.get("../zona/" + idDaZona + "/esquinas", function(esquinas) {
		var lista = []
		var idAtual = $("#esquina").val()
		for (i in esquinas) {
			var id = esquinas[i].id
			var checked = id == idAtual? "checked" : "";
			var elemento = "<label for='" + id + "'>" +
				"<input type='radio' class='selecao-esquina' name='selecao-esquina' id='" +
					esquinas[i].id + "' value='" + esquinas[i].id + "' " + checked + "/> " +
				esquinas[i].cruzamento + "<br/></label>"

			lista.push(elemento)
		}

		resultados = lista
		$("#esquinas-exibidas").html(lista)
	})
}

function atualizarEsquinas() {
	$("#zona").change(function() {
		var idDaZona = $("#zona > :selected").val();
		obterEsquinas(idDaZona)
	})
}

function registrarFiltro() {
	var buscaAnterior = ""
	resultados = $("#esquinas-exibidas").html()

	$('#filtro').keyup(function() {
		var busca = $("#filtro").val()
		var re = new RegExp(busca, "i");

		if (busca.length < buscaAnterior.length) {
			$("#esquinas-exibidas").html(resultados)
		}

		var html = $("#esquinas-exibidas").children("label").filter(function() {
			return $(this).text().search(re) != -1;
		})

		$('#esquinas-exibidas').html(html);

		buscaAnterior = busca;
	});
}
