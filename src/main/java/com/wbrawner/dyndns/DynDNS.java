package com.wbrawner.dyndns;

import com.google.cloud.functions.HttpFunction;
import com.google.cloud.functions.HttpRequest;
import com.google.cloud.functions.HttpResponse;

import java.util.List;

public class DynDNS implements HttpFunction {
    @Override
    public void service(HttpRequest request, HttpResponse response) throws Exception {
        response.setStatusCode(200);
        var writer = response.getWriter();
        for (var header : request.getHeaders().keySet()) {
            List<String> values = request.getHeaders().get(header);
            writer.write(String.format("%s: %s\n", header, String.join(", ", values)));
        }
    }
}
