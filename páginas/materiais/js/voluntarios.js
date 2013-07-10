$(function() {
	selecaoDeLider();
	selecaoDeZona();
	confirmacao();
});

function selecaoDeLider() {
	configurarBusca();

	$('#selecionar-lider').click(function(e) {
		$('#selecao-zona').slideUp(400);
		$('#selecao-lider').slideDown(400);
	});
}

function configurarBusca() {
	var comprimento = 3;
	var buscaNoServidor = "";
	var buscaAnterior = "";
	var resultados = [];

	var obterLideres = function(busca) {
		$.get("busca/lideres?contem=" + busca, function(lideres) {
			var lista = [];
			for (i in lideres) {
				var id = lideres[i].id
				var elemento = "<label for='lider" + id + "'>" +
				"<input type='radio' class='lideres' name='lideres' id='lider" + id + "'/> " +
				lideres[i].nome + "<br/></label>";

				lista.push(elemento);
			}

			$("#lideres").html(lista);
			resultados = lista;
		});
	}

	$("#busca-lider").keyup(function() {
		var novaBusca = $("#busca-lider").val();

		if (novaBusca.length == comprimento) {
			if (novaBusca != buscaNoServidor) {
				buscaNoServidor = novaBusca;
				obterLideres(buscaNoServidor);
			} else {
				$("#lideres").html(resultados);
			}

		} else if (novaBusca.length < comprimento) {
			$("#lideres").html("");

		} else {
			if (novaBusca.length < buscaAnterior.length) {
				$('#lideres').html(resultados);
			}

			var html = $("#lideres").children("label").filter(function() {
				re = new RegExp(novaBusca, "i");
				return $(this).text().search(re) != -1;
			})

			$('#lideres').html(html);
		}

		buscaAnterior = novaBusca;
	});
}

function selecaoDeZona() {
	$('#selecionar-zona').click(function(e) {
		$('#selecao-zona').slideDown(400);
		$('#selecao-lider').slideUp(400);
	});
}

function confirmacao() {
	var textoOriginal = $("#texto-turnos").text()

	$('#enviar').click(function(e) {
		var turnos = $.map($("#turnos > input:checked + label"), function(x) {
			return $(x).text()
		});

		if (turnos.length == 0) {
			$("#texto-turnos").text("Você não selecionou nenhum turno!");
		} else {
			$("#texto-turnos").text(textoOriginal)
		}
		$("#confirmacao-turnos").html(turnos.map(function(turno) {return "<li>" + turno + "</li>"}));
		$("#confirmacao-zona").html($("#zona > :selected").html());

		e.preventDefault();
		$('#confirmacao').dialog({
			resizable: false,
			dialogClass: "dialog",
			height: 330,
			width: 480,
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
